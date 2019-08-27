// Create by Yale 2018/1/8 14:47
package gorpool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Job function
type Job func()

type empty struct{}

type worker struct {
	workerPool  chan *worker
	jobQueue    chan Job
	stop        chan struct{}
	idleTime    time.Duration
	workerCount *int32
	stopWork    int32
}
type Pool struct {
	dispatcher       *dispatcher
	wg               sync.WaitGroup
	enableWaitForAll bool
}
type dispatcher struct {
	workerPool  chan *worker
	jobQueue    chan Job
	stop        chan struct{}
	idleTime    time.Duration
	workerCount int32
	workers     map[*worker]empty
}

func newWorker(workerPool chan *worker, idleTime time.Duration, workerCount *int32) *worker {

	return &worker{
		workerPool:  workerPool,
		jobQueue:    make(chan Job),
		stop:        make(chan struct{}),
		idleTime:    idleTime,
		workerCount: workerCount,
	}
}

// One worker start to work
func (w *worker) start() {
	var idTime = w.idleTime
	if idTime <= 0 {
		idTime = time.Hour * 24 * 365 * 100
	}
	tc := time.NewTimer(idTime)

	for {
		w.workerPool <- w
		select {
		case <-tc.C:
			sw := atomic.CompareAndSwapInt32(&w.stopWork, 0, 1)
			if sw {
				atomic.AddInt32(w.workerCount, -1)
				fmt.Println("stop worker success")
			} else {
				fmt.Println("stop worker fail")
			}
			return
		case job := <-w.jobQueue:

			job()
			tc = time.NewTimer(idTime)
		case <-w.stop:
			fmt.Println("worker start w.stop ")
			w.stop <- struct{}{}
			return
		}
	}

}

func (dis *dispatcher) startWorker() {

	if int(atomic.LoadInt32(&dis.workerCount)) < cap(dis.workerPool) {
		atomic.AddInt32(&dis.workerCount, 1)
		worker := newWorker(dis.workerPool, dis.idleTime, &dis.workerCount)
		go worker.start()
	}
}

// Dispatch job to free worker
func (dis *dispatcher) dispatch() {
	for {
		select {
		case job := <-dis.jobQueue:

			for {
				worker, ok := <-dis.workerPool
				if !ok {
					fmt.Println("dispatch  <-dis.workerPool ok false")
					break
				}

				if _, ok := dis.workers[worker]; !ok {
					dis.workers[worker] = empty{}
				}
				if worker.stopWork == 0 {
					worker.jobQueue <- job
					break
				}
				fmt.Println("dispatch  delete work")
				delete(dis.workers, worker)
				worker = nil
			}

		case <-dis.stop:

			for v, _ := range dis.workers {
				fmt.Println("dispatch  <-dis.stop  dis.workers")
				if v != nil && atomic.CompareAndSwapInt32(&v.stopWork, 0, 1) {
					fmt.Println("dispatch  <-dis.stop  CompareAndSwapInt32")
					v.stop <- struct{}{}
					<-v.stop
				}
			}
			dis.stop <- struct{}{}

			return
		}
	}
}
func newDispatcher(workerPool chan *worker, jobQueue chan Job) *dispatcher {
	return &dispatcher{workerPool: workerPool, jobQueue: jobQueue, stop: make(chan struct{}), workers: make(map[*worker]empty, 0)}
}

// WorkerNum is worker number of worker pool,on worker have one goroutine
//
// JobNum is job number of job pool
func NewPool(workerNum, jobNum int) *Pool {
	workers := make(chan *worker, workerNum)
	jobs := make(chan Job, jobNum)
	pool := &Pool{
		dispatcher:       newDispatcher(workers, jobs),
		enableWaitForAll: false,
	}
	return pool

}

// After idleTime stop a worker circulate, default 100 year
func (p *Pool) SetIdleDuration(idleTime time.Duration) *Pool {
	if p.dispatcher != nil {
		p.dispatcher.idleTime = idleTime
	}
	return p
}
func (p *Pool) WorkerLength() int {
	if p.dispatcher != nil {
		return int(atomic.LoadInt32(&p.dispatcher.workerCount))
	}
	return 0
}

//Add one job to job pool
func (p *Pool) AddJob(job Job) {

	p.dispatcher.startWorker()

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
	close(p.dispatcher.workerPool)
	p.dispatcher.stop <- struct{}{}
	<-p.dispatcher.stop
}

//Enable whether  WaitForAll
func (p *Pool) EnableWaitForAll(enable bool) *Pool {
	p.enableWaitForAll = enable
	return p
}

//Start worker pool and dispatch
func (p *Pool) Start() *Pool {
	go p.dispatcher.dispatch()
	return p
}
