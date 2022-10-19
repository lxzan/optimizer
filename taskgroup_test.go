package optimizer

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewLimiter(t *testing.T) {
	as := assert.New(t)
	t.Run("success", func(t *testing.T) {
		mu := sync.Mutex{}
		listA := make([]uint8, 0)
		listB := make([]uint8, 0)
		ctl := NewTaskGroup(8)
		for i := 0; i < 100; i++ {
			ctl.Push(uint8(i))
			listB = append(listB, uint8(i))
		}
		ctl.OnMessage = func(options interface{}) error {
			mu.Lock()
			listA = append(listA, options.(uint8))
			mu.Unlock()
			return nil
		}
		ctl.StartAndWait()
		as.ElementsMatch(listA, listB)
	})

	t.Run("error", func(t *testing.T) {
		ctl := NewTaskGroup(8)
		ctl.Push(1, 2, 3)
		ctl.OnMessage = func(options interface{}) error {
			return errors.New("test")
		}
		ctl.StartAndWait()
		err := ctl.Err()
		as.Error(err)
		as.Equal("test: test: test", err.Error())
	})
}
