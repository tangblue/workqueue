package workqueue_test

import (
	"context"
	"fmt"
	"github.com/tangblue/workqueue"
	"sync/atomic"
	"time"
)

type delayWork struct {
	Name  string
	Delay time.Duration
}

func (dw delayWork) Do(ctx context.Context) {
	time.Sleep(dw.Delay)
	if id := getQueueID(ctx); id >= 0 {
		fmt.Printf("%+v: %+v Delay %s seconds\n", ctx, id, dw.Delay)
	} else {
		fmt.Printf("%+v: Delay %s seconds\n", ctx, dw.Delay)
	}
}

type queueIDKey struct{}
type delayWorkContext struct {
	id int32
}

func getQueueID(ctx context.Context) int32 {
	if id, ok := ctx.Value(queueIDKey{}).(int32); ok {
		return id
	}

	return -1
}
func (dwc *delayWorkContext) Setup() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, queueIDKey{}, atomic.AddInt32(&dwc.id, 1))
	fmt.Printf("%+v: Setup\n", ctx)
	return ctx
}
func (dwc *delayWorkContext) Teardown(ctx context.Context) {
	fmt.Printf("%+v: Teardown\n", ctx)
}

func ExampleWorkQueue() {
	dwq := workqueue.New(2, 1, &delayWorkContext{})

	if err := dwq.Enqueue(delayWork{"hi", time.Second}); err != nil {
		fmt.Println("Error:", err)
	}
	dwq.Close()
	// Output:
	// context.Background.WithValue(type workqueue_test.queueIDKey, val <not Stringer>): Setup
	// context.Background.WithValue(type workqueue_test.queueIDKey, val <not Stringer>): 1 Delay 1s seconds
	// context.Background.WithValue(type workqueue_test.queueIDKey, val <not Stringer>): Teardown
}
