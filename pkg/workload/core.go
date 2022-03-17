// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package workload

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/data"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/operation"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/sample"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/selector"
	"github.com/cloudspannerecosystem/gcsb/pkg/schema"
	"github.com/cloudspannerecosystem/gcsb/pkg/workload/pool"
	"github.com/olekukonko/tablewriter"
	"github.com/rcrowley/go-metrics"
)

var (
	// Assert that WorkerPool implements Workload
	_ Workload = (*CoreWorkload)(nil)
)

const (
	defaultBufferLen = 5000
)

type (
	CoreWorkload struct {
		Context         context.Context
		Config          *config.Config
		Schema          schema.Schema
		MetricsRegistry metrics.Registry

		// Internals
		pool   *pool.PipedPool
		wg     sync.WaitGroup
		client *spanner.Client

		DataWriteGenerationTimer metrics.Timer // Used to time data generation
		DataReadGenerationTimer  metrics.Timer // Used to time data geenration
		DataWriteTimer           metrics.Timer // Used to time writes
		DataWriteMeter           metrics.Meter // Used to measure volume of writes
		DataReadTimer            metrics.Timer // Used to time reads
		DataReadMeter            metrics.Meter // Used to measure volume of reads

		// Plans and targets
		plan []*Target // The entire run plan. 1 target per table
	}
)

// NewCoreWorkload initializes a "worker pool" type workload
func NewCoreWorkload(cfg WorkloadConfig) (Workload, error) {
	wl := &CoreWorkload{
		Context:         cfg.Context,
		Config:          cfg.Config,
		Schema:          cfg.Schema,
		MetricsRegistry: cfg.MetricRegistry,
		plan:            make([]*Target, 0),
		pool: pool.NewPipedPool(pool.PipedPoolConfig{
			Workers:         cfg.Config.Threads,
			EnableOutput:    true,
			BufferOutput:    true,
			OutputBufferLen: defaultBufferLen,
			BufferInput:     true,
			InputBufferLen:  defaultBufferLen,
		}),
	}

	// Validat that metrics registry is not nil
	if wl.MetricsRegistry == nil {
		return nil, errors.New("missing metrics registry")
	}

	// Validate that schema is not nil
	if wl.Schema == nil {
		return nil, errors.New("missing schema")
	}

	// Validate that config is not nil
	if wl.Config == nil {
		return nil, errors.New("missing config")
	}

	err := wl.Initialize()
	if err != nil {
		return nil, fmt.Errorf("initializing CoreWorkload: %s", err.Error())
	}

	return wl, nil
}

func (c *CoreWorkload) Initialize() error {
	var err error
	// Initialize spanner client
	c.client, err = c.Config.Client(c.Context)
	if err != nil {
		return err
	}

	// Start the thread pool
	c.pool.Start()

	// Create our job metrics
	c.DataWriteGenerationTimer = metrics.GetOrRegisterTimer("operations.write.data", c.MetricsRegistry) // Used to time data generation
	c.DataReadGenerationTimer = metrics.GetOrRegisterTimer("operations.read.data", c.MetricsRegistry)   // Used to time data geenration
	c.DataWriteTimer = metrics.GetOrRegisterTimer("operations.write.time", c.MetricsRegistry)           // Used to time writes
	c.DataWriteMeter = metrics.GetOrRegisterMeter("operations.write.rate", c.MetricsRegistry)           // Used to measure volume of writes
	c.DataReadTimer = metrics.GetOrRegisterTimer("operations.read.time", c.MetricsRegistry)             // Used to time reads
	c.DataReadMeter = metrics.GetOrRegisterMeter("operations.read.rate", c.MetricsRegistry)             // Used to measure volume of reads

	return nil
}

// Plan will create *Targets for each TargetName
func (c *CoreWorkload) Plan(pt JobType, targets []string) error {
	var needOperationMultiplication bool
	apexTables := make([]schema.Table, 0)

	// search func for looking if targets contains the given string
	contains := func(items []string, searchterm string) bool {
		for _, item := range items {
			if item == searchterm {
				return true
			}
		}
		return false
	}

	// We can only run against one table at a time, so only expand interleaved tables if we are loading
	if pt == JobLoad {
		// Make a pass over targets and add interleaved tables that may not exist
		for _, t := range targets {
			st := c.Schema.GetTable(t)
			if st == nil {
				return fmt.Errorf("table '%s' missing from information schema", t)
			}

			// If the table is interleaved, find it's entire lineage and add it to the target list
			if st.IsInterleaved() && !st.IsApex() {
				needOperationMultiplication = true // Used below
				relatives := st.GetAllRelationNames()
				for _, n := range relatives {
					if n == t { // Avoid inserting t twice for some reason... i dont have time to figure out why this is happenign
						continue
					}
					if !contains(targets, n) {
						targets = append(targets, n)
					}
				}
			}
		}
	}

	// Iterate over targets and create Target
	for _, t := range targets {
		// Fetch table from schema by name
		st := c.Schema.GetTable(t)
		if st == nil {
			return fmt.Errorf("table '%s' missing from information schema", t)
		}

		if st.IsInterleaved() && st.IsApex() {
			apexTables = append(apexTables, st) // Collect a slice of apex tables
		}

		// Create target
		target := &Target{
			Config:                   c.Config,
			Context:                  c.Context,
			Client:                   c.client,
			JobType:                  pt,
			Table:                    st,
			TableName:                t,
			ColumnNames:              st.ColumnNames(),
			DataWriteGenerationTimer: c.DataWriteGenerationTimer,
			DataReadGenerationTimer:  c.DataReadGenerationTimer,
			DataWriteTimer:           c.DataWriteTimer,
			DataWriteMeter:           c.DataWriteMeter,
			DataReadTimer:            c.DataReadTimer,
			DataReadMeter:            c.DataReadMeter,
		}

		// If we are in 'run' context
		if pt == JobRun {
			// Generate an operation selector
			sel, err := c.GetOperationSelector()
			if err != nil {
				return fmt.Errorf("creating operation selector: %s", err.Error())
			}

			target.OperationSelector = sel

			// If our read fraction is > 0,
			// We have faith that the operation selector will not return reads if read fraction is <= 0
			if c.Config.Operations.Read > 0 {
				// Sample the table and create a sample generator
				sg, err := c.GetReadGeneratorMap(target.Table)
				if err != nil {
					return fmt.Errorf("creating sample generator: %s", err.Error())
				}

				target.ReadGenerator = sg
			}
		}

		// Create a generator map for the table
		gm, err := c.GetGeneratorMap(target.Table)
		if err != nil {
			return fmt.Errorf("creating generator map: %s", err.Error())
		}

		target.WriteGenerator = gm

		////
		// Try to set an operations count on the target
		////

		// See if this table is in the config
		ct := c.Config.Table(target.TableName)
		if ct == nil {
			// There is no configuration for this table
			if target.Table.IsInterleaved() {
				if target.Table.IsApex() {
					target.Operations = c.Config.Operations.Total
				} else {
					target.Operations = config.DefaultTableOperations
				}
			} else {
				target.Operations = c.Config.Operations.Total
			}
		} else {
			// Table is in the configuration but has no operations config
			if ct.Operations == nil {
				if target.Table.IsInterleaved() {
					if target.Table.IsApex() {
						target.Operations = c.Config.Operations.Total // if it's the apex table apply total operations
					} else {
						target.Operations = config.DefaultTableOperations // A default table operations multiplier for child tables
					}
				} else {
					target.Operations = c.Config.Operations.Total
				}
			} else {
				// Table does have a configuration value for operations in config file
				target.Operations = ct.Operations.Total
			}
		}

		c.plan = append(c.plan, target)
	}

	// So if our phase is load, the operations per target are actually multipliers. Now we go through and do that multiplication
	if pt == JobLoad && needOperationMultiplication {
		for _, at := range apexTables {
			apexTarget := FindTargetByName(c.plan, at.Name())

			lastOps := apexTarget.Operations
			relatives := at.GetAllRelationNames()
			for _, cn := range relatives {
				if cn == at.Name() {
					continue
				}

				relativeTarget := FindTargetByName(c.plan, cn)
				relativeTarget.Operations = lastOps * relativeTarget.Operations
				lastOps = relativeTarget.Operations
			}
		}
	}

	return nil
}

func (c *CoreWorkload) Load(x []string) error {
	// Plan our run
	err := c.Plan(JobLoad, x)
	if err != nil {
		return fmt.Errorf("planning run: %s", err.Error())
	}

	// Summarize plan
	c.SummarizePlan()

	// Execute our run
	err = c.Execute()
	if err != nil {
		return fmt.Errorf("executing run: %s", err.Error())
	}

	return nil
}

// Run will execute a Run phase against the target table.
func (c *CoreWorkload) Run(x string) error {
	// Fetch table from schema
	table := c.Schema.GetTable(x)
	if table == nil {
		return fmt.Errorf("table '%s' missing from schema", x)
	}

	// Check if table is interleaved
	if table.IsInterleaved() {
		// Check if table is apex. If not, return error
		if !table.IsApex() {
			apex := table.GetApex()
			return fmt.Errorf("can only execute run against apex table (try '%s')", apex.Name())
		}
	}

	// Plan our run
	err := c.Plan(JobRun, []string{x})
	if err != nil {
		return fmt.Errorf("planning run: %s", err.Error())
	}

	// Summarize plan
	c.SummarizePlan()

	// Execute our run
	err = c.Execute()
	if err != nil {
		return fmt.Errorf("executing run: %s", err.Error())
	}

	return nil
}

func (c *CoreWorkload) Execute() error {
	////
	// Setup transition threads
	////

	var abortErr error           // If we abort for some reason, we will assign the reason to this error and return it
	var timeout <-chan time.Time // The timeout channel. It is never used if max execution time is not set
	abort := make(chan struct{}) // Used to halt everything on fatal error
	done := make(chan struct{})  // Signaled when the waitgroup is done (all jobs exit normally)

	// If max execution time is set and is > 0, setup a timer that will fire on the timeout chan
	if c.Config.MaxExecutionTime > 0 {
		to := time.NewTimer(c.Config.MaxExecutionTime)
		timeout = to.C
	}

	// Create a waitgroup thread. This thread listens to the output of c.pool and decrements
	// the wait group when the job is complete
	waitGroupChan := make(chan pool.Job, defaultBufferLen)
	waitGroupEnd := make(chan bool)
	waitGroupFunc := func() {
		for {
			select {
			case <-waitGroupEnd:
				return
			case j := <-waitGroupChan:
				job, ok := j.(*Job)
				if !ok {
					panic("received job that was not *Job (BUG)")
				}

				// If the job has a fatal error, we will abort
				if job.FatalErr != nil {
					abortErr = job.FatalErr
					abort <- struct{}{}
					close(abort)
					return
				}

				c.wg.Done() // Must release! Otherwise we will deadlock
			}
		}
	}

	c.pool.BindPool(waitGroupChan)
	go waitGroupFunc()

	// Cleanup on return
	// defer func() {
	// 	c.pool.Stop()
	// 	waitGroupEnd <- true
	// }()

	////
	// Do work. Generate jobs and feed them to the pool
	////
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		for _, target := range c.plan {
			// Bucketize operations
			buckets := c.bucketOps(target.Operations, c.Config.Threads)

			// For each bucket of operations, make a job
			for _, ops := range buckets {
				// Get a job from the target
				job := target.NewJob()

				// Set operations
				job.Operations = ops

				// Submit job to pool
				c.pool.Submit(job)

				// Increment waitgorup
				c.wg.Add(1)
			}
		}
	}()

	////
	// Wait for jobs to process
	////
	go func() {
		// Wait for all jobs to flow through the pipeline
		c.wg.Wait()
		done <- struct{}{}
		close(done)
	}()

	select {
	case <-done: // all jobs exited normally
		return nil
	case <-abort: // a job encountered a fatal error
		return abortErr
	case <-timeout: // max execution time reached
		return errors.New("max execution time reached")
	}
}

func (c *CoreWorkload) Stop() error {
	if c.pool != nil {
		c.pool.Stop()
	}

	return nil
}

// GetReadGeneratorMap will sample rows from the table and create a map structure for creating point reads
func (c *CoreWorkload) GetReadGeneratorMap(t schema.Table) (*sample.SampleGenerator, error) {
	samples, err := c.SampleTable(t)
	if err != nil {
		return nil, fmt.Errorf("sampling table: %s", err.Error())
	}

	return generator.GetReadGeneratorMap(samples, t.PrimaryKeyNames())
}

// SampleTable will return a map[string]interface of values using the tables primary keys
func (c *CoreWorkload) SampleTable(t schema.Table) (map[string]interface{}, error) {
	return generator.SampleTable(c.Config, c.Context, c.client, t)
}

// GetGeneratorMap will return a generator map suitable for creating insert operations against a table
func (c *CoreWorkload) GetGeneratorMap(t schema.Table) (data.GeneratorMap, error) {
	return generator.GetDataGeneratorMapForTable(*c.Config, t)
}

func (c *CoreWorkload) GetOperationSelector() (selector.Selector, error) {
	return operation.NewOperationSelector(c.Config)
}

// bucketOps will divide operations into buckets and grow each bucket to handle remainders
func (c *CoreWorkload) bucketOps(n int, k int) []int {
	r := make([]int, k)

	e := n / k
	o := n % k
	for i := 0; i <= k-1; i++ {
		if o > 0 {
			r[i] = e + 1
			o--
		} else {
			r[i] = e
		}
	}

	return r
}

func (c *CoreWorkload) SummarizePlan() {
	tableString := &strings.Builder{}
	t := tablewriter.NewWriter(tableString)
	t.SetHeader([]string{
		"Table", "Operations", "Read", "Write", "Context",
	})

	for _, target := range c.plan {
		l := []string{
			target.TableName,
			fmt.Sprintf("%d", target.Operations),
		}

		if target.JobType == JobRun {
			l = append(l,
				fmt.Sprintf("%d", c.Config.Operations.Read),
				fmt.Sprintf("%d", c.Config.Operations.Write),
			)
		} else {
			l = append(l, "N/A", "N/A")
		}

		if target.JobType == JobLoad {
			l = append(l, "LOAD")
		}

		if target.JobType == JobRun {
			l = append(l, "RUN")
		}

		t.Append(l)
	}

	t.Render()
	scanner := bufio.NewScanner(strings.NewReader(tableString.String()))
	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}
