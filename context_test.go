package workqueue_test

import (
	"fmt"
	"github.com/tangblue/workqueue"
	"testing"
	"time"
)

func TestContextNil(t *testing.T) {
	dwq := workqueue.NewWorkQueue(4, 2, nil)

	for i := 0; i < 5; i += 1 {
		if err := dwq.QueueWork(delayWork{"hi", time.Second}); err != nil {
			fmt.Println("Error:", err)
		}
	}
	dwq.Stop()
}

func TestWorkQueue(t *testing.T) {
	dwq := workqueue.NewWorkQueue(4, 2, delayWorkContext{})

	for i := 0; i < 5; i += 1 {
		if err := dwq.QueueWork(delayWork{"hi", time.Second}); err != nil {
			fmt.Println("Error:", err)
		}
	}
	dwq.Stop()
}
