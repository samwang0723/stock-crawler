// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package concurrent

import (
	"context"
	"sync"

	log "github.com/samwang0723/stock-crawler/internal/logger"
)

type Dispatcher struct {
	// a pool of worker channels that registered with dispatcher
	workerPool chan JobChan
	workers    []*Worker
	maxWorkers int
	waitGroup  *sync.WaitGroup
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan JobChan, maxWorkers)
	return &Dispatcher{
		workerPool: pool,
		maxWorkers: maxWorkers,
		waitGroup:  &sync.WaitGroup{},
	}
}

func (d *Dispatcher) Run(ctx context.Context) {
	d.waitGroup.Add(d.maxWorkers)
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i, d.workerPool, d.waitGroup)
		worker.Start()
		d.workers = append(d.workers, worker)
	}

	go d.dispatch(ctx)
}

func (d *Dispatcher) dispatch(ctx context.Context) {
	for {
		select {
		// job request received
		case job, ok := <-JobQueue:
			if ok {
				// try to obtain a available worker job channel
				// this will block until a worker is idle
				jobChannel := <-d.workerPool

				// dispatch the job to worker job channel
				jobChannel <- job
			}
		case <-ctx.Done():
			log.Warn("!!! dispatch: context cancelled !!!")
			for _, w := range d.workers {
				w.Stop()
			}
			return
		}
	}
}

func (d *Dispatcher) WaitGroup() *sync.WaitGroup {
	return d.waitGroup
}
