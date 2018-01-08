# gorpool

Simple Goroutine pool 

# How to use

```
go get github.com/yale8848/gorpool
```

## Simple example

```go

package main

import (
	"github.com/yale8848/gorpool"
	"time"
	"fmt"
)

func main() {

	p := gorpool.NewPool(5, 10).Start()
	defer p.StopAll()
	for i := 0; i < 100; i++ {
		count := i
		p.AddJob(func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("%d\r\n", count)
		})

	}
	time.Sleep(2 * time.Second)
}

```

## Wait

```go

package main

import (
	"fmt"
	"github.com/yale8848/gorpool"
	"time"
)

func main() {

	p := gorpool.NewPool(5, 10).
		Start().
		EnableWaitForAll(true)
	for i := 0; i < 100; i++ {
		count := i
		p.AddJob(func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("%d\r\n", count)
		})
	}
	p.WaitForAll()
	p.StopAll()
}

```


