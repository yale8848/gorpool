// Create by Yale 2018/1/8 14:48
package gorpool

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	p := NewPool(5, 10).Start()
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


func TestPool_EnableWaitForAll(t *testing.T) {
	p := NewPool(5, 10).Start().
		EnableWaitForAll(true)
	for i := 0; i < 100; i++ {
		count := i
		p.AddJob(func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf(" %d\r\n",count)
		})
	}
	p.WaitForAll()
	p.StopAll()
}

func TestPool_StopAll(t *testing.T) {
	rnum := runtime.NumGoroutine()
	p := NewPool(5, 10).Start().
		EnableWaitForAll(true)
	defer func() {
		p.StopAll()
		if rnum != runtime.NumGoroutine() {
			t.Error("Goroutine not stop")
		} else {
			t.Log("Goroutine  stoped")
		}
	}()
	for i := 0; i < 100; i++ {
		count := i
		p.AddJob(func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("%d\r\n", count)
		})
	}
	p.WaitForAll()
}
