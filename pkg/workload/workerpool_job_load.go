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

	"cloud.google.com/go/spanner"
	"github.com/rcrowley/go-metrics"
	"github.com/cloudspannerecosystem/gcsb/pkg/generator/data"
	"github.com/cloudspannerecosystem/gcsb/pkg/workload/pool"
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
		Batch           bool
		BatchSize       int
		WaitGroup       *sync.WaitGroup
		MetricsRegistry metrics.Registry
	}
)

func (j *WorkerPoolLoadJob) Execute() {
	if j.Batch {
		j.InsertMapBatch()
	} else {
		j.InsertMap()
	}

	j.WaitGroup.Done()
}

func (j *WorkerPoolLoadJob) InsertMapBatch() {
	batch := make([]*spanner.Mutation, 0, j.BatchSize)

	for i := 1; i <= j.RowCount; i++ {
		m := make(map[string]interface{}, len(j.GeneratorMap))
		for k, v := range j.GeneratorMap {
			m[k] = v.Next()
		}

		batch = append(batch, spanner.InsertMap(j.TableName, m))

		if len(batch) == j.BatchSize {
			_, err := j.Client.Apply(j.Context, batch)
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

			batch = nil
			batch = make([]*spanner.Mutation, 0, j.BatchSize)
		}
	}

	// Flush the buffer at the end
	if len(batch) > 0 {
		_, err := j.Client.Apply(j.Context, batch)
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
}

func (j *WorkerPoolLoadJob) InsertMap() {
	for i := 0; i <= j.RowCount; i++ {
		m := make(map[string]interface{}, len(j.GeneratorMap))
		for k, v := range j.GeneratorMap {
			m[k] = v.Next()
		}

		_, err := j.Client.Apply(j.Context, []*spanner.Mutation{spanner.InsertMap(j.TableName, m)})
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
}

func (j *WorkerPoolLoadJob) InsertDML() {
	for i := 0; i <= j.RowCount; i++ {
		stmt := spanner.NewStatement(j.Statement)
		for k, v := range j.GeneratorMap {
			stmt.Params[k] = v.Next()
		}

		_, err := j.Client.ReadWriteTransaction(j.Context, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
			_ = txn.Query(ctx, stmt)
			return nil
		})

		if err != nil {
			log.Printf("error in write transaction: %s", err.Error())
		}
	}
}
