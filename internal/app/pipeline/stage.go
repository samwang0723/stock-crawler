// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package pipeline

import (
	"context"
	"sync"
	"time"

	"github.com/samwang0723/stock-crawler/internal/retry"
	"golang.org/x/xerrors"
)

const (
	defaultRetryTimes = 3
)

type fifo struct {
	proc Processor
}

// FIFO returns a StageRunner that processes incoming payloads in a first-in
// first-out fashion. Each input is passed to the specified processor and its
// output is emitted to the next stage.
func FIFO(proc Processor) StageRunner {
	return fifo{proc: proc}
}

// Run implements StageRunner.
func (r fifo) Run(ctx context.Context, params StageParams) {
	for {
		select {
		case <-ctx.Done():
			// Asked to cleanly shut down
			return
		case payloadIn, ok := <-params.Input():
			if !ok {
				return
			}

			payloadOut, err := r.proc.Process(ctx, payloadIn)
			if err != nil {
				wrappedErr := xerrors.Errorf("pipeline stage %d: %w", params.StageIndex(), err)
				maybeEmitError(wrappedErr, params.Error())
			}
			// If the processor did not output a payload for the
			// next stage there is nothing we need to do.
			if payloadOut == nil {
				payloadIn.MarkAsProcessed()

				continue
			}

			// Output processed data
			select {
			case params.Output() <- payloadOut:
			case <-ctx.Done():
				// Asked to cleanly shut down
				return
			}
		}
	}
}

type dynamicWorkerPool struct {
	proc      Processor
	tokenPool chan struct{}
	interval  time.Duration
}

// DynamicWorkerPool returns a StageRunner that maintains a dynamic worker pool
// that can scale up to maxWorkers for processing incoming inputs in parallel
// and emitting their outputs to the next stage.
func DynamicWorkerPool(proc Processor, maxWorkers int, interval time.Duration) StageRunner {
	if maxWorkers <= 0 {
		panic("DynamicWorkerPool: maxWorkers must be > 0")
	}

	tokenPool := make(chan struct{}, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		tokenPool <- struct{}{}
	}

	return &dynamicWorkerPool{proc: proc, tokenPool: tokenPool, interval: interval}
}

//nolint:nolintlint, cyclop // Run implements StageRunner.
func (p *dynamicWorkerPool) Run(ctx context.Context, params StageParams) {
stop:
	for {
		select {
		case <-ctx.Done():
			// Asked to cleanly shut down
			break stop
		case payloadIn, ok := <-params.Input():
			if !ok {
				break stop
			}

			var token struct{}
			select {
			case token = <-p.tokenPool:
			case <-ctx.Done():
				break stop
			}

			go func(payloadIn Payload, token struct{}) {
				defer func() { p.tokenPool <- token }()

				var payloadOut Payload

				err := retry.Retry(defaultRetryTimes, p.interval, func() error {
					out, procErr := p.proc.Process(ctx, payloadIn)
					payloadOut = out

					if procErr != nil {
						return xerrors.Errorf("retry error: %w", procErr)
					}

					return nil
				})
				if err != nil {
					wrappedErr := xerrors.Errorf("pipeline stage %d: %w", params.StageIndex(), err)
					maybeEmitError(wrappedErr, params.Error())

					return
				}

				// If the processor did not output a payload for the
				// next stage there is nothing we need to do.
				if payloadOut == nil {
					payloadIn.MarkAsProcessed()

					return
				}

				// Output processed data
				select {
				case params.Output() <- payloadOut:
				case <-ctx.Done():
				}
			}(payloadIn, token)

			// prevent rate limit
			<-time.After(p.interval)
		}
	}

	// Wait for all workers to exit by trying to empty the token pool
	for i := 0; i < cap(p.tokenPool); i++ {
		<-p.tokenPool
	}
}

type broadcast struct {
	fifos []StageRunner
}

// Broadcast returns a StageRunner that passes a copy of each incoming payload
// to all specified processors and emits their outputs to the next stage.
func Broadcast(procs ...Processor) StageRunner {
	if len(procs) == 0 {
		panic("Broadcast: at least one processor must be specified")
	}

	fifos := make([]StageRunner, len(procs))
	for i, p := range procs {
		fifos[i] = FIFO(p)
	}

	return &broadcast{fifos: fifos}
}

// Run implements StageRunner.
//
//nolint:nolintlint, cyclop
func (b *broadcast) Run(ctx context.Context, params StageParams) {
	var (
		waitGroup sync.WaitGroup
		inCh      = make([]chan Payload, len(b.fifos))
	)

	// Start each FIFO in a go-routine. Each FIFO gets its own dedicated
	// input channel and the shared output channel passed to Run.
	for index := 0; index < len(b.fifos); index++ {
		waitGroup.Add(1)

		inCh[index] = make(chan Payload)

		go func(fifoIndex int) {
			fifoParams := &workerParams{
				stage: params.StageIndex(),
				inCh:  inCh[fifoIndex],
				outCh: params.Output(),
				errCh: params.Error(),
			}
			b.fifos[fifoIndex].Run(ctx, fifoParams)
			waitGroup.Done()
		}(index)
	}

done:
	for {
		// Read incoming payloads and pass them to each FIFO
		select {
		case <-ctx.Done():
			break done
		case payload, ok := <-params.Input():
			if !ok {
				break done
			}
			for index := len(b.fifos) - 1; index >= 0; index-- {
				// As each FIFO might modify the payload, to
				// avoid data races we need to make a copy of
				// the payload for all FIFOs except the first.
				fifoPayload := payload
				if index != 0 {
					fifoPayload = payload.Clone()
				}
				select {
				case <-ctx.Done():
					break done
				case inCh[index] <- fifoPayload:
					// payload sent to i_th FIFO
				}
			}
		}
	}

	// Close input channels and wait for FIFOs to exit
	for _, ch := range inCh {
		close(ch)
	}

	waitGroup.Wait()
}
