// Create by Yale 2018/10/11 17:41
package gpool

type Worker interface {
	AddJob( job Job)
}
type worker struct {

	jobChan chan Job
	stopChan chan bool
	pool pool
}
func (w *worker)AddJob(job Job){
	w.jobChan<-job
}
func (w *worker)work()  {

	for{
		select {
		    case job:= <- w.jobChan:
				w.pool.addWait()
		    	job()
		    	w.pool.waitDone()
		    	w.pool.putWorker(w)
		    case stop:=<-w.stopChan:
				if stop {
					return
				}
		}
	}


}
