// Create by Yale 2018/1/8 14:47
package gorpool

import (
	"sync"
)

type Job func()
type worker struct {
	workerPool chan *worker
	jobQueue   chan Job
	stop       chan struct{}
}
type Pool struct {
	dispatcher       *dispatcher
	wg               sync.WaitGroup
	enableWaitForAll bool
}
type dispatcher struct {
	workerPool chan *worker
	jobQueue   chan Job
	stop       chan struct{}
}

func newWorker(workerPool chan *worker) *worker {

	return &worker{
		workerPool: workerPool,
		jobQueue:   make(chan Job),
		stop:       make(chan struct{}),
	}
}
func (w *worker) start() {
	for {
		w.workerPool <- w
		select {
		case job := <-w.jobQueue:
			job()
		case <-w.stop:
			w.stop <- struct{}{}
			return
		}

	}

}

func (dis *dispatcher) dispatch() {
	for {
		select {
		case job := <-dis.jobQueue:
			worker := <-dis.workerPool
			worker.jobQueue <- job
		case <-dis.stop:
			for i := 0; i < cap(dis.workerPool); i++ {
				worker := <-dis.workerPool
				worker.stop <- struct{}{}
				<-worker.stop
			}
			dis.stop <- struct{}{}
			return
		}
	}
}
func newDispatcher(workerPool chan *worker, jobQueue chan Job) *dispatcher {
	return &dispatcher{workerPool: workerPool, jobQueue: jobQueue, stop: make(chan struct{})}
}

func NewPool(workerNum, jobNum int) *Pool {
	workers := make(chan *worker, workerNum)
	jobs := make(chan Job, jobNum)

	pool := &Pool{
		dispatcher:       newDispatcher(workers, jobs),
		enableWaitForAll: false,
	}

	return pool

}
func (p *Pool) AddJob(job Job) {
	if p.enableWaitForAll {
		p.wg.Add(1)
	}

	p.dispatcher.jobQueue <- func() {
		job()
		if p.enableWaitForAll {
			p.wg.Done()
		}
	}
}

func (p *Pool) WaitForAll() {
	if p.enableWaitForAll {
		p.wg.Wait()
	}
}

func (p *Pool) StopAll() {
	p.dispatcher.stop <- struct{}{}
	<-p.dispatcher.stop
}
func (p *Pool) EnableWaitForAll(enable bool) *Pool {
	p.enableWaitForAll = enable
	return p
}
func (p *Pool) Start() *Pool {

	for i := 0; i < cap(p.dispatcher.workerPool); i++ {
		worker := newWorker(p.dispatcher.workerPool)
		go worker.start()
	}
	go p.dispatcher.dispatch()

	return p
}
