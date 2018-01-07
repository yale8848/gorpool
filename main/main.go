package main

import (
	"gopool"
	"fmt"
	"time"
)

func main()  {

	p:=gopool.NewPool(1,1)

	p.Start()

	for i:=0;i<100 ;i++  {
		count := i
		p.AddJob(func() {
			fmt.Printf("%d\r\n",count)
		})

	}
	time.Sleep(2*time.Second)
}