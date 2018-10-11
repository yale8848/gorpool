// Create by Yale 2018/10/11 18:39
package gpool

import (
	"testing"
	"fmt"
)

func TestDefaultPool(t *testing.T) {


	cht:=make(chan int,3)
	fmt.Println(len(cht))
	fmt.Println(cap(cht))

	cht<-1

	fmt.Println(len(cht))
	fmt.Println(cap(cht))
	cht<-1

	fmt.Println(len(cht))
	fmt.Println(cap(cht))

	cht<-1

	fmt.Println(len(cht))
	fmt.Println(cap(cht))

/*	p:=DefaultPool()

	for i:=0;i<10;i++{
		it:=i
		p.GetWorker().AddJob(func() {
			fmt.Println(it)
		})
	}
	p.Wait()*/



}