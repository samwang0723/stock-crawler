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
	"sync"
	"time"

	log "github.com/samwang0723/stock-crawler/internal/logger"
	"github.com/samwang0723/stock-crawler/internal/retry"
)

type Job interface {
	Do() error
}

// define job channel
type JobChan chan Job

type Worker struct {
	workerPool chan JobChan
	// job channel is for worker to get job
	jobChannel JobChan
	quit       chan bool
	waitGroup  *sync.WaitGroup
	id         int
}

var (
	// buffered channel to send worker requests on
	JobQueue JobChan
)

func NewWorker(id int, workerPool chan JobChan, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:         id,
		workerPool: workerPool,
		jobChannel: make(JobChan),
		quit:       make(chan bool),
		waitGroup:  wg,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			// keep register available job channel back to worker pool
			w.workerPool <- w.jobChannel
			select {
			case job := <-w.jobChannel:
				if err := retry.Retry(3, 2000*time.Millisecond, job.Do); err != nil {
					log.Errorf("worker(%d) job execution with failure: %+v", w.id, err)
				}
			// received quit event and terminate worker
			case <-w.quit:
				log.Warnf("worker(%d) terminated: context cancelled!", w.id)
				w.waitGroup.Done()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.quit <- true
}
