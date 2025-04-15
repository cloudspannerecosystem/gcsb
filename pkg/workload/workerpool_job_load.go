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
	"sync"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/data"
	"github.com/cloudspannerecosystem/gcsb/pkg/workload/pool"
	"github.com/rcrowley/go-metrics"
	"google.golang.org/grpc/codes"
)

var (
	// Assert that WorkerPoolLoadJob implements pool.Job
	_ pool.Job = (*WorkerPoolLoadJob)(nil)
)

type (
	// WorkerPoolLoadJob is responsible for inserting data into a table
	WorkerPoolLoadJob struct {
		Context         context.Context
		Client          *spanner.Client
		TableName       string
		RowCount        int
		Statement       string
		GeneratorMap    data.GeneratorMap
		CommitDelay     time.Duration // Maximum commit delay for throughput optimization
		WaitGroup       *sync.WaitGroup
		MetricsRegistry metrics.Registry
	}
)

func (j *WorkerPoolLoadJob) Execute() {
	j.InsertMapBatch()
	j.WaitGroup.Done()
}

func (j *WorkerPoolLoadJob) InsertMapBatch() {
	// Log the commit delay value before starting transaction
	log.Printf("Worker pool starting transaction with commit delay: %v", j.CommitDelay)

	// Create a transaction with commit options
	commitDelay := j.CommitDelay // Create a local variable to get its address
	_, err := j.Client.ReadWriteTransactionWithOptions(j.Context,
		func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			for i := 1; i <= j.RowCount; i++ {
				m := make(map[string]interface{}, len(j.GeneratorMap))
				for k, v := range j.GeneratorMap {
					m[k] = v.Next()
				}

				// Buffer the mutation in the transaction
				err := txn.BufferWrite([]*spanner.Mutation{
					spanner.InsertMap(j.TableName, m),
				})
				if err != nil {
					return err
				}
			}
			return nil
		},
		spanner.TransactionOptions{
			CommitOptions: spanner.CommitOptions{
				MaxCommitDelay:    &commitDelay,
				ReturnCommitStats: true,
			},
		},
	)

	// Log the commit delay value being used
	log.Printf("Worker pool transaction completed with commit delay: %v", commitDelay)

	if err != nil {
		sErr := spanner.ErrCode(err)
		if sErr == codes.Canceled {
			return
		}
		if sErr == codes.Unauthenticated {
			log.Println("Received unrecoverable authentication error. Worker is exiting.")
			return
		}
		log.Printf("error in write transaction: %s", err.Error())
	}
}
