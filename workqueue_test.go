package workqueue

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
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

func TestContextNil(t *testing.T) {
	dwq := NewWorkQueue(4, 2, nil)

	for i := 0; i < 5; i += 1 {
		if err := dwq.QueueWork(delayWork{"hi", time.Second}); err != nil {
			fmt.Println("Error:", err)
		}
	}
	dwq.Stop()
}

func TestContext(t *testing.T) {
	dwq := NewWorkQueue(4, 2, delayWorkContext{})

	for i := 0; i < 5; i += 1 {
		if err := dwq.QueueWork(delayWork{"hi", time.Second}); err != nil {
			fmt.Println("Error:", err)
		}
	}
	dwq.Stop()
}
