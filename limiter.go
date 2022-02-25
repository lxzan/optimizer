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

func (this *Limiter) Push(eles ...interface{}) {
	for _, ele := range eles {
		this.q.Push(ele)
		if atomic.LoadInt64(&this.curNum) < this.maxNum {
			this.do()
		}
	}
}

func (this *Limiter) do() {
	var item = this.q.Front()
	if item == nil {
		return
	}

	atomic.AddInt64(&this.curNum, 1)
	go func(doc interface{}) {
		this.handler(doc)
		atomic.AddInt64(&this.curNum, -1)
		this.do()
	}(item)
}

func (this *Limiter) Stop() {
	docs := this.q.Clear()
	for _, item := range docs {
		this.handler(item)
	}
}
