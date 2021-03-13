package optimizer

import (
	"fmt"
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	var t0 = time.Now().UnixNano()
	r := NewLimiter(10, func(doc interface{}) {
		var t1 = time.Now().UnixNano()
		fmt.Printf("idx=%d, time=%dms\n", doc.(int), (t1-t0)/1000000)
	})
	r.Start()

	for i := 0; i < 20; i++ {
		r.Push(i)
	}

	time.Sleep(3 * time.Second)
}