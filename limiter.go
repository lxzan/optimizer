package optimizer

import (
	"context"
	"sync/atomic"
	"time"
)

type Limiter struct {
	*Queue
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
		Queue:    NewQueue(),
		maxNum:   num,
		curNum:   0,
		interval: time.Duration(interval) * time.Millisecond,
		ctx:      ctx,
		cancel:   cancel,
		handler:  handler,
	}
	return o
}

func (c *Limiter) Start() {
	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if atomic.LoadInt64(&c.curNum) >= c.maxNum {
					continue
				}
				if doc := c.Front(); doc != nil {
					atomic.AddInt64(&c.curNum, 1)
					go func() {
						c.handler(doc)
						atomic.AddInt64(&c.curNum, -1)
					}()
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *Limiter) Stop() {
	c.cancel()
	docs := c.Clear()
	for _, item := range docs {
		c.handler(item)
	}
}
