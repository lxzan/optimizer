package optimizer

import (
	"context"
	"time"
)

type BatchProcessor struct {
	*Queue
	interval  time.Duration
	ctx       context.Context
	cancel    context.CancelFunc
	OnMessage func(docs []interface{})
}

func NewBatchProcessor(ctx context.Context, interval time.Duration) *BatchProcessor {
	cctx, cancel := context.WithCancel(ctx)
	return &BatchProcessor{
		Queue:    NewQueue(),
		interval: interval,
		ctx:      cctx,
		cancel:   cancel,
	}
}

func (c *BatchProcessor) Start() {
	go func() {
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if arr := c.Clear(); len(arr) > 0 {
					c.OnMessage(arr)
				}
			case <-c.ctx.Done():
				return
			}
		}
	}()
}

func (c *BatchProcessor) Stop() {
	c.cancel()
	if arr := c.Clear(); len(arr) > 0 {
		c.OnMessage(arr)
	}
}
