package workerpool

type Worker struct {
	ID      int
	jobChan chan Job
	pool    chan chan Job
	quit    chan bool
}

func NewWorker(ID int, pool chan chan Job) *Worker {
	return &Worker{
		ID:      ID,
		jobChan: make(chan Job),
		pool:    pool,

		quit: make(chan bool),
	}
}

func (wr *Worker) Run() {
	for {
		wr.pool <- wr.jobChan
		select {
		case job := <-wr.jobChan:
			// fmt.Printf("🎬 worker received job%d \n", job.Id)
			job.Run()
		case <-wr.quit:
			return

		}
	}

}

func (wr *Worker) Stop() {
	wr.quit <- true
}
