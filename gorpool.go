// Create by Yale 2018/1/8 14:47
package gorpool

import (
	"sync"
	"sync/atomic"
	"time"
)

const MaxTime = time.Hour * 24 * 365 * 100

// Job function
type Job func()

type worker struct {
	workerPool chan *worker
	jobQueue   chan Job
	stop       chan struct{}
	idleTime   time.Duration

	dis    *dispatcher
	isStop bool
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

	locker sync.Mutex
}

func newWorker(workerPool chan *worker, idleTime time.Duration, dis *dispatcher) *worker {

	return &worker{
		workerPool: workerPool,
		jobQueue:   make(chan Job),
		stop:       make(chan struct{}, 1),
		idleTime:   idleTime,
		dis:        dis,
	}
}

func (dis *dispatcher) startWorker() {

	dis.locker.Lock()
	defer dis.locker.Unlock()

	if dis.workerCount < int32(cap(dis.workerPool)) {
		dis.workerCount++
		worker := newWorker(dis.workerPool, dis.idleTime, dis)
		go worker.start()
	}
}

// One worker start to work
func (w *worker) start() {

	var idTime = w.idleTime
	if idTime <= 0 {
		idTime = MaxTime
	}
	tc := time.NewTimer(idTime)

	for {
		w.workerPool <- w
		select {
		case <-tc.C:

			w.dis.locker.Lock()
			if !w.isStop {
				w.isStop = true
				w.dis.workerCount--
			} else {
				<-w.stop
			}
			w.dis.locker.Unlock()

			return
		case job := <-w.jobQueue:
			job()
			if idTime != MaxTime {
				tc = time.NewTimer(idTime)
			}
		case <-w.stop:
			return
		}
	}

}

// Dispatch job to free worker
func (dis *dispatcher) dispatch() {
	for {
		select {
		case job := <-dis.jobQueue:

			worker := <-dis.workerPool
			if !worker.isStop {
				worker.jobQueue <- job
			} else {
				dis.startWorker()
				dis.jobQueue <- job
			}

		case <-dis.stop:

			func() {
				dis.locker.Lock()
				defer dis.locker.Unlock()

				wl := len(dis.workerPool)
				for i := 0; i < wl; i++ {
					worker := <-dis.workerPool
					if !worker.isStop {
						worker.isStop = true
						worker.stop <- struct{}{}
					}

				}
			}()
			time.Sleep(100 * time.Millisecond)
			dis.stop <- struct{}{}
			return
		}
	}
}
func newDispatcher(workerPool chan *worker, jobQueue chan Job) *dispatcher {
	return &dispatcher{workerPool: workerPool, jobQueue: jobQueue, stop: make(chan struct{})}
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

// Add one job to job pool
func (p *Pool) AddJob(job Job) {

	p.dispatcher.startWorker()

	p.dispatcher.jobQueue <- func() {
		if p.enableWaitForAll {
			p.wg.Add(1)
		}
		job()
		if p.enableWaitForAll {
			p.wg.Done()
		}
	}

}

// Wait all job finish
func (p *Pool) WaitForAll() {
	if p.enableWaitForAll {
		p.wg.Wait()
	}
}

// Stop all worker
func (p *Pool) StopAll() {
	p.dispatcher.stop <- struct{}{}
	<-p.dispatcher.stop
}

// Enable whether  WaitForAll
func (p *Pool) EnableWaitForAll(enable bool) *Pool {
	p.enableWaitForAll = enable
	return p
}

// Start worker pool and dispatch
func (p *Pool) Start() *Pool {
	go p.dispatcher.dispatch()
	return p
}
