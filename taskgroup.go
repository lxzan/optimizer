package optimizer

import (
	"sync"
	"sync/atomic"
)

type TaskGroup struct {
	collector *errorCollector // error collector
	q         *Queue          // task queue
	goroutine int64           // max concurrent coroutine
	taskDone  int64           // completed tasks
	taskTotal int64           // total tasks
	handler   func(ctl *TaskGroup, options interface{}) error
	signal    chan bool
}

// goroutine: max concurrent coroutine
func NewTaskGroup(goroutine int64, handler func(ctl *TaskGroup, options interface{}) error) *TaskGroup {
	if goroutine <= 0 {
		goroutine = 8
	}
	o := &TaskGroup{
		collector: &errorCollector{mu: &sync.RWMutex{}},
		q:         NewQueue(),
		goroutine: goroutine,
		taskDone:  0,
		signal:    make(chan bool),
		handler:   handler,
	}
	return o
}

func (c *TaskGroup) Len() int {
	return c.q.Len()
}

func (c *TaskGroup) Push(eles ...interface{}) {
	atomic.AddInt64(&c.taskTotal, int64(len(eles)))
	c.q.Push(eles...)
}

func (c *TaskGroup) do() {
	if atomic.LoadInt64(&c.taskDone) == atomic.LoadInt64(&c.taskTotal) {
		c.signal <- true
		return
	}

	if item := c.q.Front(); item != nil {
		go func(doc interface{}) {
			if err := c.handler(c, doc); err != nil {
				c.collector.MarkFailedWithError(err)
			} else {
				c.collector.MarkSucceed()
			}
			atomic.AddInt64(&c.taskDone, 1)
			c.do()
		}(item)
	}
}

func (c *TaskGroup) StartAndWait() {
	var length = atomic.LoadInt64(&c.taskTotal)
	if length == 0 {
		go func() {
			c.signal <- true
		}()
		return
	}

	var co = min(int(c.goroutine), int(length))
	for i := 0; i < co; i++ {
		c.do()
	}

	<-c.signal
}

func (c *TaskGroup) Err() error {
	return c.collector.Err()
}
