package workerpool

type Worker struct {
	ID      int
	jobChan <-chan *Job

	quit chan bool
}

func NewWorker(channel <-chan *Job, ID int) *Worker {
	return &Worker{
		ID:      ID,
		jobChan: channel,
		quit:    make(chan bool),
	}
}

func (wr *Worker) Run() {
	go func() {
		for {
			select {
			case job := <-wr.jobChan:
				// fmt.Printf("🎬 worker received job%d \n", job.Id)
				job.Run()
			case <-wr.quit:
				return

			}
		}
	}()

}

func (wr *Worker) Stop() {
	go func() {
		wr.quit <- true
	}()
}
