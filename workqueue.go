package workqueue

import (
	"context"
	"errors"
	"sync"
)

// Work is the interface that wraps the work for the work queue.
type Work interface {
	// Do does the real work in context of a worker.
	Do(context.Context)
}

type WorkContextHelper interface {
	Setup() context.Context
	Teardown(context.Context)
}

type WorkQueue struct {
	workQueue chan Work
	wch       WorkContextHelper
	wg        *sync.WaitGroup
}

func (wq *WorkQueue) QueueWork(w Work) error {
	select {
	case wq.workQueue <- w:
	default:
		return errors.New("Work queue full")
	}

	return nil
}

func (wq *WorkQueue) Stop() {
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

func NewWorkQueue(nworks, nworkers int, wch WorkContextHelper) *WorkQueue {
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
