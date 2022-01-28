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
	"context"
	"errors"
	"fmt"
	"sync"

	"cloud.google.com/go/spanner"
	"github.com/rcrowley/go-metrics"
	"github.com/cloudspannerecosystem/gcsb/pkg/config"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/operation"
	"github.com/cloudspannerecosystem/gcsb/pkg/schema"
	"github.com/cloudspannerecosystem/gcsb/pkg/workload/pool"
)

var (
	// Assert that WorkerPool implements Workload
	_ Workload = (*WorkerPool)(nil)
)

type (
	WorkerPool struct { // Implement Workload
		Context         context.Context
		Config          *config.Config
		Schema          schema.Schema
		initialized     bool
		Pool            *pool.Pool
		Jobs            []pool.Job
		wg              sync.WaitGroup
		client          *spanner.Client
		MetricsRegistry metrics.Registry
	}
)

// NewPoolWorkload initializes a "worker pool" type workload
func NewPoolWorkload(cfg WorkloadConfig) (Workload, error) {
	w := &WorkerPool{
		Context:         cfg.Context,
		Config:          cfg.Config,
		Schema:          cfg.Schema,
		Jobs:            make([]pool.Job, 0),
		MetricsRegistry: cfg.MetricRegistry,
		Pool: pool.NewPool(pool.PoolConfig{
			Workers:        cfg.Config.Threads,
			BufferInput:    true,
			InputBufferLen: 100, // TODO: don't hardcode this
		}),
	}

	return w, nil
}

func (w *WorkerPool) Initialize() error {
	if w.MetricsRegistry == nil {
		return errors.New("missing metrics registry")
	}

	var err error
	w.client, err = w.Config.Client(w.Context)
	if err != nil {
		return err
	}

	w.Pool.Start()

	w.initialized = true

	return nil
}

// Stop the worker pool
func (w *WorkerPool) Stop() error {
	if w.Pool != nil {
		w.Pool.Stop()
	}

	return nil
}

func (w *WorkerPool) Load(tables []string) error {
	if len(tables) <= 0 {
		return fmt.Errorf("need 1 table")
	}

	tableName := tables[0]

	if !w.initialized {
		err := w.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize workload: %s", err.Error())
		}

	}

	// Construct generator map for table
	table := w.Schema.GetTable(tableName)
	if table == nil {
		return fmt.Errorf("table '%s' missing from schema", tableName)
	}

	opsPerJob := w.Config.Operations.Total / w.Config.Threads
	for i := 1; i <= w.Config.Threads; i++ {
		// Create a unique generator map instance for each job
		genMap, err := generator.GetDataGeneratorMapForTable(*w.Config, table)
		if err != nil {
			return fmt.Errorf("getting generator map: %s", err.Error())
		}

		// For fun lets grab an insert statement just in case we decide to use dml later
		stmt, err := table.PointInsertStatement()
		if err != nil {
			return fmt.Errorf("getting table write statement: %s", err.Error())
		}

		j := &WorkerPoolLoadJob{
			Context:         w.Context,
			Client:          w.client,
			TableName:       tableName,
			RowCount:        opsPerJob,
			Statement:       stmt,
			GeneratorMap:    genMap,
			Batch:           true,
			BatchSize:       5,
			WaitGroup:       &w.wg,
			MetricsRegistry: w.MetricsRegistry,
		}

		w.Jobs = append(w.Jobs, j)
		w.wg.Add(1)
		w.Pool.Submit(j)
	}

	w.wg.Wait()

	return nil
}

func (w *WorkerPool) Run(tableName string) error {
	// Initialize the pool
	if !w.initialized {
		err := w.Initialize()
		if err != nil {
			return fmt.Errorf("failed to initialize workload: %s", err.Error())
		}
	}

	// Grab table from schema
	table := w.Schema.GetTable(tableName)
	if table == nil {
		return fmt.Errorf("table '%s' missing from schema", tableName)
	}

	// Determine number of operations per thread
	opsPerJob := w.Config.Operations.Total / w.Config.Threads

	// Need to fetch primary key(s) from target table
	pKeyNames := table.PrimaryKeyNames()
	if len(pKeyNames) <= 0 {
		return errors.New("unable to determine primary key column for READ operations")
	}

	// Use primary keys to TABLESAMPLE
	samples, err := generator.SampleTable(w.Config, w.Context, w.client, table)
	if err != nil {
		return fmt.Errorf("error sampling table: %s", err.Error())
	}

	// Create 1 job per thread
	for i := 1; i <= w.Config.Threads; i++ {
		// Create operation selector
		sel, err := operation.NewOperationSelector(w.Config)
		if err != nil {
			return fmt.Errorf("getting operation selector: %s", err.Error())
		}

		// Construct generator map for table inserts
		insertMap, err := generator.GetDataGeneratorMapForTable(*w.Config, table)
		if err != nil {
			return fmt.Errorf("getting insert generator map: %s", err.Error())
		}

		// initialize a static value generator for READ ops (readMap)
		gen, err := generator.GetReadGeneratorMap(samples, table.PrimaryKeyNames())
		if err != nil {
			return fmt.Errorf("error getting read generator: %s", err.Error())
		}

		// Create Job
		j := &WorkerPoolRunJob{
			Context:           w.Context,
			Client:            w.client,
			TableName:         tableName,
			ReadGenerator:     gen,
			WriteMap:          insertMap,
			OperationSelector: sel,
			WaitGroup:         &w.wg,
			StaleReads:        w.Config.Operations.ReadStale,
			Staleness:         w.Config.Operations.Staleness,
			Operations:        opsPerJob,
			Table:             table,
			MetricsRegistry:   w.MetricsRegistry,
		}

		// Keep a reference to the job for no reason
		w.Jobs = append(w.Jobs, j)

		// Increment waitgroup
		w.wg.Add(1)

		// Submit job
		w.Pool.Submit(j)
	}

	// Block until all jobs complete
	w.wg.Wait()

	return nil
}
