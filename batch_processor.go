package optimizer

import (
	"context"
	"time"
)

type BatchProcessor struct {
	*Queue
	interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	handler  func([]interface{})
}

func NewBatchProcessor(interval time.Duration, handler func(docs []interface{})) *BatchProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	return &BatchProcessor{
		Queue:    NewQueue(),
		interval: interval,
		ctx:      ctx,
		cancel:   cancel,
		handler:  handler,
	}
}

func (this *BatchProcessor) Start() {
	go func() {
		ticker := time.NewTicker(this.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if arr := this.Clear(); len(arr) > 0 {
					this.handler(arr)
				}
			case <-this.ctx.Done():
				return
			}
		}
	}()
}

func (this *BatchProcessor) Stop() {
	this.cancel()
	if arr := this.Clear(); len(arr) > 0 {
		this.handler(arr)
	}
}
