package optimizer

import (
	"fmt"
	"testing"
	"time"
)

func TestNewBatchProcessor(t *testing.T) {
	var p = NewBatchProcessor(500*time.Millisecond, func(docs []interface{}) {
		for _, item := range docs {
			println(item.(int))
		}
		fmt.Print("\n")
	})
	p.Push(1, 3, 5)
	p.Start()

	time.Sleep(600 * time.Millisecond)
	p.Push(7, 9)
	p.Stop()
	p.Push(11, 13)
}
