// Create by Yale 2018/1/8 14:48
package gorpool

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	rn := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := NewPool(3, 10).SetIdleDuration(1 * time.Second).Start()
	defer p.StopAll()
	for i := 0; i < 10; i++ {
		count := i
		p.AddJob(func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("v %d\r\n", count)
		})
		p.AddJob(func() {
			time.Sleep(1 * time.Second)
			fmt.Printf("v %d\r\n", count+100)
		})
		time.Sleep(time.Duration(rn.Intn(3)) * time.Second)
		fmt.Printf("w %d\r\n", p.WorkerLength())
	}
	time.Sleep(60 * time.Second)
}

func TestPool_EnableWaitForAll(t *testing.T) {
	p := NewPool(5, 10).Start().
		EnableWaitForAll(true)
	for i := 0; i < 100; i++ {
		count := i
		p.AddJob(func() {
			time.Sleep(10 * time.Millisecond)
			fmt.Printf(" %d\r\n", count)
		})
	}
	p.WaitForAll()
	p.StopAll()
}

func TestPool_StopAll(t *testing.T) {
	rn := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnum := runtime.NumGoroutine()
	p := NewPool(5, 10).SetIdleDuration(300 * time.Millisecond).Start().
		EnableWaitForAll(true)
	defer func() {
		p.StopAll()

		time.Sleep(5 * time.Second)
		if rnum != runtime.NumGoroutine() {
			t.Error("Goroutine not stop")
		} else {
			t.Log("Goroutine  stoped")
		}
	}()
	for i := 0; i < 100; i++ {
		count := i
		p.AddJob(func() {
			time.Sleep(time.Duration(rn.Intn(300)) * time.Millisecond)
			fmt.Printf("%d\r\n", count)
		})
		time.Sleep(time.Duration(rn.Intn(600)) * time.Millisecond)
	}
	p.WaitForAll()
}
