package optimizer

import (
	"sync/atomic"
)

type Limiter struct {
	q       *Queue // task queue
	maxNum  int64  // maxNum: max concurrent coroutine
	curNum  int64  // curNum: current concurrent coroutine
	handler func(opt interface{})
}

// maxNum: max concurrent coroutine
func NewLimiter(maxNum int64, handler func(doc interface{})) *Limiter {
	o := &Limiter{
		q:       NewQueue(),
		maxNum:  maxNum,
		curNum:  0,
		handler: handler,
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
	var item = c.q.Front()
	if item == nil {
		return
	}

	atomic.AddInt64(&c.curNum, 1)
	go func(doc interface{}) {
		c.handler(doc)
		atomic.AddInt64(&c.curNum, -1)
		c.do()
	}(item)
}

func (c *Limiter) Stop() {
	docs := c.q.Clear()
	for _, item := range docs {
		c.handler(item)
	}
}
