// Package workqueue provides basic interfaces to work queue.
package workqueue

import (
	"context"
	"errors"
	"sync"
)

// ErrDropped is the error returned when the work queue is full.
var ErrDropped = errors.New("Dropped")

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
	*sync.WaitGroup
}

// Enqueue queues the work to the work queue. The queued work will be done in a worker's context.
func (wq *WorkQueue) Enqueue(w Work) error {
	select {
	case wq.workQueue <- w:
	default:
		return ErrDropped
	}

	return nil
}

// Close closes the channel of the work queue and wait all queued works done.
func (wq *WorkQueue) Close() {
	close(wq.workQueue)
	wq.Wait()
}

func (wq *WorkQueue) doWorks() {
	var ctx context.Context
	if wq.wch != nil {
		ctx = wq.wch.Setup()
		defer wq.wch.Teardown(ctx)
	} else {
		ctx = context.Background()
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
		WaitGroup: &sync.WaitGroup{},
	}

	wq.Add(nworkers)
	for i := 0; i < nworkers; i++ {
		go func() {
			wq.doWorks()
			wq.Done()
		}()
	}

	return &wq
}
