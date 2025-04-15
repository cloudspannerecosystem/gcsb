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
	"log"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/data"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/operation"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/sample"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/selector"
	"github.com/rcrowley/go-metrics"
	"google.golang.org/grpc/codes"
)

type (
	Job struct {
		JobType           JobType           // Job Type (load or run)
		Context           context.Context   // Context
		Client            *spanner.Client   // Spanner Client
		Table             string            // Table name to execute against
		Operations        int               // How many operations in this job
		CommitDelay       time.Duration     // Maximum commit delay for throughput optimization
		Columns           []string          // Tables column names to ask for during reads
		StaleReads        bool              // Perform stale reads if true
		Staleness         time.Duration     // If performing stale reads, use this exact staleness
		OperationSelector selector.Selector // Weghted choice selector (read or write)

		// Generators
		WriteGenerator data.GeneratorMap       // Generator for making row data
		ReadGenerator  *sample.SampleGenerator // Generator for point reads

		// Metrics
		DataWriteGenerationTimer metrics.Timer // Used to time data generation
		DataReadGenerationTimer  metrics.Timer // Used to time data geenration
		DataWriteTimer           metrics.Timer // Used to time writes
		DataWriteMeter           metrics.Meter // Used to measure volume of writes
		DataReadTimer            metrics.Timer // Used to time reads
		DataReadMeter            metrics.Meter // Used to measure volume of reads

		FatalErr error
	}

	// A simplified transaction interface to consolidate stale vs strong reads
	transaction interface {
		ReadRow(ctx context.Context, table string, key spanner.Key, columns []string) (*spanner.Row, error)
		Close()
	}
)

const (
	// For batched inserts, how many rows per API call
	DefaultBatchSize = 100
)

func (j *Job) Execute() {
	switch j.JobType {
	case JobLoad: // Load data to table
		// Insert all operations in a single transaction with commit delay
		err := j.InsertWithCommitDelay()
		if err != nil { // If err is returned, it is fatal
			return
		}
	case JobRun: // Run against table
		// Generate $operations reads/writes
		for i := 0; i <= j.Operations; i++ {
			// Select an operation to perform
			op := j.OperationSelector.Select().Item().(operation.Operation)
			switch op {
			case operation.READ:
				err := j.ReadOne()
				if err != nil { // If err is returned, it is fatal
					return
				}
			case operation.WRITE:
				err := j.InsertOne()
				if err != nil { // if err is returned, it is fatal
					return
				}
			}
		}
	default:
		log.Printf("unknown JobType(%d)", j.JobType)
		return
	}
}

func (j *Job) ReadOne() error {
	// Log the commit delay value before starting transaction
	log.Printf("Benchmark read transaction starting with commit delay: %v", j.CommitDelay)

	// Generate read predicate
	r := j.generateReadKey()

	// Get a read transaction
	tx := j.getReadTransaction()

	// perform read
	err := j.readRow(tx, r)

	// Log the commit delay value being used
	log.Printf("Benchmark read transaction completed with commit delay: %v", j.CommitDelay)

	// Check for fatal errors
	err = j.checkSpannerError(err)

	return err
}

/*
 * InsertOne will insert one row into the jobs table
 */
func (j *Job) InsertOne() error {
	// Log the commit delay value before starting transaction
	log.Printf("Benchmark transaction starting with commit delay: %v", j.CommitDelay)

	// Create a map of row data
	m := j.generateRow()

	// Insert the row using the mutation API with commit delay
	_, err := j.Client.Apply(j.Context,
		[]*spanner.Mutation{
			spanner.InsertMap(j.Table, m),
		},
		spanner.ApplyCommitOptions(spanner.CommitOptions{
			MaxCommitDelay:    &j.CommitDelay,
			ReturnCommitStats: true,
		}),
	)

	// Log the commit delay value being used
	log.Printf("Benchmark transaction completed with commit delay: %v", j.CommitDelay)

	// Check if error is fatal. Since we're only performing 1 op,
	// This is called mostly for it's side effect of collecting
	// non-fatal errors to be used elsewhere
	_ = j.checkSpannerError(err)

	return err
}

func (j *Job) InsertWithCommitDelay() error {
	// Log the commit delay value before starting transaction
	log.Printf("Starting transaction with commit delay: %v", j.CommitDelay)

	// Create a transaction with commit options
	_, err := j.Client.ReadWriteTransactionWithOptions(j.Context,
		func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			for i := 1; i <= j.Operations; i++ {
				// Generate a map for the row data
				m := j.generateRow()

				// Buffer the mutation in the transaction
				err := txn.BufferWrite([]*spanner.Mutation{
					spanner.InsertMap(j.Table, m),
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
		spanner.TransactionOptions{
			CommitOptions: spanner.CommitOptions{
				MaxCommitDelay:    &j.CommitDelay,
				ReturnCommitStats: true,
			},
		},
	)

	// Log the commit delay value being used
	log.Printf("Transaction completed with commit delay: %v", j.CommitDelay)

	// Check if error is fatal
	err = j.checkSpannerError(err)
	return err
}

// checkSpannerError will return the error if it is fatal,
// if not, it will collect the error and return nil
func (j *Job) checkSpannerError(err error) error {
	if err != nil {
		spannerErr := spanner.ErrCode(err)

		// If error is codes.Unauthenticated, return. We can not proceed
		if spannerErr == codes.Unauthenticated {
			j.FatalErr = err // Set the FatalError so we halt the entire workload
			return err
		}

		// If error is codes.Canceled, return. Our context is canceled
		if spannerErr == codes.Canceled {
			j.FatalErr = err // Set the FatalError so we halt the entire workload
			return err
		}

		// TODO: Collect errors
	}

	return nil
}

// generateRow will return a map of row data based on the jobs GeneratorMap
func (j *Job) generateRow() map[string]interface{} {
	// Generate a map for the row data
	m := make(map[string]interface{}, len(j.WriteGenerator))
	j.DataWriteGenerationTimer.Time(func() {
		for k, v := range j.WriteGenerator {
			m[k] = v.Next()
		}
	})

	return m
}

// generateReadKey will return a spanner.Key suitable for executing a point read
func (j *Job) generateReadKey() spanner.Key {
	var r spanner.Key
	j.DataReadGenerationTimer.Time(func() {
		r = j.ReadGenerator.Next().(spanner.Key)
	})

	return r
}

// getTransactions will return a read transactions based on jobs config
func (j *Job) getReadTransaction() transaction {
	if j.StaleReads {
		return j.Client.ReadOnlyTransaction().WithTimestampBound(spanner.ExactStaleness(j.Staleness))
	}

	return j.Client.Single()
}

// applyMutations will call apply on a slice of spanner mutations and return any errors
func (j *Job) applyMutations(muts []*spanner.Mutation) error {
	var err error
	j.DataWriteTimer.Time(func() {
		_, err = j.Client.Apply(j.Context, muts)
		j.DataWriteMeter.Mark(int64(len(muts))) // Mark how many write mutations were proccessed
	})

	return err
}

// readRow will query the table for the provided spanner.Key using the passed transaction
func (j *Job) readRow(tx transaction, r spanner.Key) error {
	var err error
	j.DataReadTimer.Time(func() {
		// Perform read, discard row
		_, err = tx.ReadRow(j.Context, j.Table, r, j.Columns)
		j.DataReadMeter.Mark(1) // measure read rate
	})

	// Close the transaction
	tx.Close()

	return err
}
