package workqueue_test

import (
	"context"
	"fmt"
	"github.com/tangblue/workqueue"
	"sync/atomic"
	"time"
)

type key string

const (
	keyID = key("id")
)

var id int32

func SetID(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyID, atomic.AddInt32(&id, 1))
}
func GetID(ctx context.Context) int32 {
	return ctx.Value(keyID).(int32)
}

type delayWork struct {
	Name  string
	Delay time.Duration
}

func (dw delayWork) Do(ctx context.Context) {
	time.Sleep(dw.Delay)
	if ctx != nil {
		fmt.Printf("%+v: %+v Delay %s seconds\n", ctx, GetID(ctx), dw.Delay)
	} else {
		fmt.Printf("%+v: Delay %s seconds\n", ctx, dw.Delay)
	}
}

type delayWorkContext struct{}

func (dwc delayWorkContext) Setup() context.Context {
	ctx := SetID(context.Background())
	fmt.Printf("%+v: Setup\n", ctx)
	return ctx
}
func (dwc delayWorkContext) Teardown(ctx context.Context) {
	fmt.Printf("%+v: Teardown\n", ctx)
}

func ExampleWorkQueue() {
	atomic.StoreInt32(&id, 0)
	dwq := workqueue.NewWorkQueue(2, 1, delayWorkContext{})

	if err := dwq.QueueWork(delayWork{"hi", time.Second}); err != nil {
		fmt.Println("Error:", err)
	}
	dwq.Stop()
	// Output:
	// context.Background.WithValue("id", 1): Setup
	// context.Background.WithValue("id", 1): 1 Delay 1s seconds
	// context.Background.WithValue("id", 1): Teardown
}
