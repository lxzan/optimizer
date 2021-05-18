package optimizer

import (
	"context"
	"sync/atomic"
	"time"
)

type Limiter struct {
	q        *Queue
	maxNum   int64
	curNum   int64
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	handler  func(opt interface{})
}

// num: max concurrent number per second, <=1000
// interval: interval of checking new task
func NewLimiter(num int64, handler func(doc interface{})) *Limiter {
	ctx, cancel := context.WithCancel(context.Background())
	interval := 1000 / num
	o := &Limiter{
		q:        NewQueue(),
		maxNum:   num,
		curNum:   0,
		interval: time.Duration(interval) * time.Millisecond,
		ctx:      ctx,
		cancel:   cancel,
		handler:  handler,
	}
	return o
}

func (c *Limiter) Push(eles ...interface{}) {
	for _, ele := range eles {
		c.q.Push(ele)
		if atomic.LoadInt64(&c.curNum) < c.maxNum {
			c.do()
		}
	}
}

func (c *Limiter) do() {
	if doc := c.q.Front(); doc != nil {
		atomic.AddInt64(&c.curNum, 1)
		go func() {
			c.handler(doc)
			atomic.AddInt64(&c.curNum, -1)
			c.do()
		}()
	}
}

func (c *Limiter) Stop() {
	c.cancel()
	docs := c.q.Clear()
	for _, item := range docs {
		c.handler(item)
	}
}
