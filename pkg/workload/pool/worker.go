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

package pool

type (
	// Worker is a goroutine that listens for work and processes incoming requests (Jobs)
	Worker struct {
		ID            int
		WorkerChannel chan chan Job
		JobOutput     chan Job // Job Output Channel
		JobInput      chan Job // Job input channel
		End           chan bool
	}
)

// Start worker
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerChannel <- w.JobInput
			select {
			case job := <-w.JobInput:
				job.Execute()
			case <-w.End:
				return
			}
		}
	}()
}

// Stop worker
func (w *Worker) Stop() {
	w.End <- true
}
