package optimizer

import "sync"

type Queue struct {
	sync.Mutex
	data []interface{}
}

func NewQueue() *Queue {
	return &Queue{
		data: make([]interface{}, 0),
	}
}

func (this *Queue) Push(eles ...interface{}) {
	this.Lock()
	this.data = append(this.data, eles...)
	this.Unlock()
}

func (this *Queue) Front() interface{} {
	this.Lock()
	defer this.Unlock()

	if n := len(this.data); n == 0 {
		return nil
	} else {
		var result = this.data[0]
		this.data = this.data[1:]
		return result
	}
}

func (this *Queue) Len() int {
	this.Lock()
	length := len(this.data)
	this.Unlock()
	return length
}

// clear queue, return all data
func (this *Queue) Clear() []interface{} {
	this.Lock()
	data := this.data
	this.data = make([]interface{}, 0)
	this.Unlock()
	return data
}
