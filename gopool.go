package gopool

type Job func()
type worker struct{
	workerPool chan *worker
	jobQueue chan Job
}
type Pool struct{
	JobQueue chan Job
	dispatcher *dispatcher
}
type dispatcher struct {
	workerPool chan *worker
	jobQueue chan Job
}

func newWorker(workerPool chan *worker) *worker  {

	return &worker{
		workerPool:workerPool,
		jobQueue:make(chan  Job),
	}
}
func (w *worker)start()   {

	go func() {
		var job Job
		for{
             w.workerPool <-w

			select {
			  case job=<-w.jobQueue:
			  	  job()
			}

		}

	}()
}

func (dis *dispatcher) dispatch()  {
	go func() {
		for{
			select {
			case worker:=<-dis.workerPool:
				job:=<-dis.jobQueue
				worker.jobQueue<-job
			}
		}
	}()
}
func newDispatcher(workerPool chan *worker,jobQueue chan Job)  *dispatcher {
	return &dispatcher{workerPool:workerPool,jobQueue:jobQueue}
}


func NewPool(workerNum ,jobNum int) *Pool  {
	workers:=make(chan *worker,workerNum)
	jobs:=make(chan Job,jobNum)

	pool:= &Pool{
		JobQueue:jobs,
		dispatcher:newDispatcher(workers,jobs),
	}

	return pool

}
func (p *Pool)AddJob(job Job )  {
	p.JobQueue<-job
}

func (p *Pool)Start( )  {


	for i:=0;i<cap(p.dispatcher.workerPool) ;  i++{
		worker:=newWorker(p.dispatcher.workerPool)
		worker.start()
	}

	p.dispatcher.dispatch()
}