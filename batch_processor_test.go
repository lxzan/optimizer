package optimizer

import (
	"testing"
	"time"
)

func TestNewBatchProcessor(t *testing.T) {
	var p = NewBatchProcessor(500*time.Millisecond, func(docs []interface{}) {
		for _, item := range docs {
			println(item.(int))
		}
	})
	p.Push(1, 3, 5)
	p.Start()
	time.Sleep(600 * time.Millisecond)
	p.Push(7)
	p.Stop()
}
