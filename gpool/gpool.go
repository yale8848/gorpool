// Create by Yale 2018/10/11 17:41
package gpool

import (
	"sync"
)

type GPool interface {

	GetWorker() Worker
	Wait()
}

type Job func()

type pool struct {
	capacity int
	wg sync.WaitGroup
	isWait bool
	workerPool chan Worker
}

func (p*pool)Wait()  {
	if p.isWait {
		p.wg.Wait()
	}
}
func DefaultPool() GPool  {
	cap:=100
	return &pool{capacity:cap,wg:sync.WaitGroup{},workerPool:make(chan Worker,cap)}
}
func (p *pool)addWait(){
	if p.isWait {
		p.wg.Add(1)
	}
}
func (p *pool)waitDone(){
	if p.isWait {
		p.wg.Done()
	}
}
func (p *pool)putWorker(w Worker){
	p.workerPool<-w
}
func (p *pool)GetWorker() Worker  {
	if len(p.workerPool) == 0 {

	}
	if len(p.workerPool) < p.capacity {
		wk:= &worker{jobChan:make(chan Job),stopChan:make(chan bool)}
		go wk.work()
		return wk

	}else{

	}

	return nil

}
