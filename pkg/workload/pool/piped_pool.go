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
	// BindFunc is called to bind the output of one pool to the input of another
	BindFunc func(Job, chan Job) error

	// PipedPoolConfig configures the pool
	PipedPoolConfig struct {
		Workers         int
		EnableOutput    bool
		BufferOutput    bool
		OutputBufferLen int
		BufferInput     bool
		InputBufferLen  int
	}

	// PipedPool will start a pool of workers and dispatch jobs to them
	PipedPool struct {
		Config        PipedPoolConfig
		WorkInput     chan Job  // For client to send work to
		WorkOutput    chan Job  // For workers to send jobs to after completion
		BinderEnd     chan bool // Used to notify binder the pool is stopping
		End           chan bool
		Workers       []PipedWorker
		isBound       bool
		WorkerChannel chan chan Job
	}
)

// NewPipedPool will return a configured PipedPool. You must call Start() to
// begin processing jobs
func NewPipedPool(cfg PipedPoolConfig) *PipedPool {

	wp := &PipedPool{
		Config:        cfg,
		End:           make(chan bool),     // channel to spin down workers
		BinderEnd:     make(chan bool),     // channel to inform binders the pool is exiting
		WorkerChannel: make(chan chan Job), // WorkerChannel is a channel of work worker channels (lol)
	}

	var input, output chan Job
	if cfg.BufferInput {
		input = make(chan Job, cfg.InputBufferLen) // channel to recieve work
	} else {
		input = make(chan Job)
	}
	wp.WorkInput = input

	if cfg.EnableOutput {
		if cfg.BufferOutput {
			output = make(chan Job, cfg.OutputBufferLen)
		} else {
			output = make(chan Job)
		}

		wp.WorkOutput = output
	}

	var i int
	for i < cfg.Workers {
		i++
		worker := PipedWorker{
			ID:            i,
			JobInput:      make(chan Job),
			WorkerChannel: wp.WorkerChannel,
			End:           make(chan bool),
		}

		if cfg.EnableOutput {
			worker.WriteToOutChannel = true
			worker.JobOutput = output
		}
		worker.Start()
		wp.Workers = append(wp.Workers, worker) // store worker
	}

	return wp
}

// Start will start the job dispatcher
func (wp *PipedPool) Start() {
	// start collector
	go func() {
		for {
			select {
			case <-wp.End:
				for _, w := range wp.Workers {
					w.Stop() // stop worker
				}
				return
			case work := <-wp.WorkInput:
				worker := <-wp.WorkerChannel // wait for available channel
				worker <- work               // dispatch work to worker
			}
		}
	}()
}

// BindPool will bind the output of this pool with the input channel of another
func (wp *PipedPool) BindPool(inputChannel chan Job) {
	wp.isBound = true
	go func() {
		for {
			select {
			case <-wp.BinderEnd:
				wp.isBound = false
				return
			case j := <-wp.WorkOutput:
				// fmt.Println("Got job on output")
				inputChannel <- j
			}
		}
	}()
}

// Submit is syntax sugar for pushing a job into the input channel
func (wp *PipedPool) Submit(job Job) {
	wp.WorkInput <- job
}

// Stop workers and then the PipedPool dispatcher
func (wp *PipedPool) Stop() {
	wp.End <- true
	if wp.isBound {
		wp.BinderEnd <- true
	}
}
