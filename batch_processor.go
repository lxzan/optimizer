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

func (c *BatchProcessor) Start() {
	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.handler(c.Clear())
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *BatchProcessor) Stop() {
	c.cancel()
	c.handler(c.Clear())
}
