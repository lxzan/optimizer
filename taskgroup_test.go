package optimizer

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewLimiter(t *testing.T) {
	as := assert.New(t)
	t.Run("", func(t *testing.T) {
		mu := sync.Mutex{}
		listA := make([]uint8, 0)
		listB := make([]uint8, 0)
		ctl := NewTaskGroup(8, func(ctl *TaskGroup, options interface{}) error {
			mu.Lock()
			listA = append(listA, options.(uint8))
			mu.Unlock()
			return nil
		})
		for i := 0; i < 100; i++ {
			ctl.Push(uint8(i))
			listB = append(listB, uint8(i))
		}
		ctl.StartAndWait()
		as.ElementsMatch(listA, listB)
	})
}
