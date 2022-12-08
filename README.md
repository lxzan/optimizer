## Optimizer

#### WorkerQueue 

> 工作队列, 可以不断往里面添加不同种类的任务, 一旦有资源空闲就去执行

```go
package main

import (
	"context"
	"fmt"
	"github.com/lxzan/optimizer"
	"time"
)

func Add(args interface{}) error {
	arr := args.([]int)
	ans := 0
	for _, item := range arr {
		ans += item
	}
	fmt.Printf("args=%v, ans=%d\n", args, ans)
	return nil
}

func Mul(args interface{}) error {
	arr := args.([]int)
	ans := 1
	for _, item := range arr {
		ans *= item
	}
	fmt.Printf("args=%v, ans=%d\n", args, ans)
	return nil
}

func main() {
	args1 := []int{1, 3}
	args2 := []int{1, 3, 5}
	w := optimizer.NewWorkerQueue(context.Background(), 8)
	w.Push(
		optimizer.Job{Args: args1, Do: Add},
		optimizer.Job{Args: args1, Do: Mul},
		optimizer.Job{Args: args2, Do: Add},
		optimizer.Job{Args: args2, Do: Mul},
	)
	w.StopAndWait(50*time.Millisecond, 30*time.Second)
}
```

```
args=[1 3], ans=4
args=[1 3 5], ans=15
args=[1 3], ans=3
args=[1 3 5], ans=9
```


#### WorkerGroup 

> 工作组, 添加一组任务, 等待任务完全被执行

```go
package main

import (
	"context"
	"fmt"
	"github.com/lxzan/optimizer"
	"sync/atomic"
)

func main() {
	sum := int64(0)
	w := optimizer.NewWorkerGroup(context.Background(), 8)
	for i := int64(1); i <= 100; i++ {
		w.Push(i)
	}
	w.OnMessage = func(options interface{}) error {
		atomic.AddInt64(&sum, options.(int64))
		return nil
	}
	w.StartAndWait()
	fmt.Printf("sum=%d\n", sum)
}
```

```
sum=5050
```