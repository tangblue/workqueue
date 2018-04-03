// Package workqueue provides basic interfaces to work queue.
package workqueue

import (
	"context"
	"errors"
	"sync"
)

// Work is the interface that wraps the Do method for the work queue.
type Work interface {
	// Do does the real work in the context of a worker.
	Do(context.Context)
}

// WorkContextHelper is the interface that wraps the work context methods.
type WorkContextHelper interface {
	// Setup setups the worker's context.
	Setup() context.Context
	// Teardown cleanups the worker's context.
	Teardown(context.Context)
}

type WorkQueue struct {
	workQueue chan Work
	wch       WorkContextHelper
	wg        *sync.WaitGroup
}

// Enqueue queues the work to the work queue. The queued work will be done in a worker's context.
func (wq *WorkQueue) Enqueue(w Work) error {
	select {
	case wq.workQueue <- w:
	default:
		return errors.New("Work queue full")
	}

	return nil
}

// Close closes the channel of the work queue and wait all queued works done.
func (wq *WorkQueue) Close() {
	close(wq.workQueue)
	wq.wg.Wait()
}

func (wq *WorkQueue) doWorks() {
	var ctx context.Context
	if wq.wch != nil {
		ctx = wq.wch.Setup()
		defer wq.wch.Teardown(ctx)
	}

	for w := range wq.workQueue {
		w.Do(ctx)
	}
}

// New creates a new work queue.
func New(nworks, nworkers int, wch WorkContextHelper) *WorkQueue {
	wq := WorkQueue{
		workQueue: make(chan Work, nworks),
		wch:       wch,
		wg:        &sync.WaitGroup{},
	}

	for i := 0; i < nworkers; i++ {
		wq.wg.Add(1)
		go func() {
			wq.doWorks()
			wq.wg.Done()
		}()
	}

	return &wq
}
