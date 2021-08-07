package workerpool

type Dispatcher struct {
	concurrency   int
	pool          chan chan Job
	jobQuequeChan chan Job
	runBackground chan bool
	workers       []*Worker
}

func NewDispatcher(concurrency int) *Dispatcher {
	return &Dispatcher{
		concurrency:   concurrency,
		runBackground: make(chan bool),
		pool:          make(chan chan Job, concurrency),
		jobQuequeChan: make(chan Job, concurrency),
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.concurrency; i++ {
		worker := NewWorker(i, d.pool)
		d.workers = append(d.workers, worker)
		go worker.Run()
	}

	go d.dispatch()

	<-d.runBackground
}

func (d *Dispatcher) dispatch() {
	for job := range d.jobQuequeChan {
		workerChan := <-d.pool // wait for available channel
		workerChan <- job      // dispatch work to worker
	}
}

func (d *Dispatcher) Enqueue(job Job) {
	d.jobQuequeChan <- job
}

func (d *Dispatcher) Stop() {
	for _, w := range d.workers {
		go w.Stop()
	}

	d.runBackground <- true
}
