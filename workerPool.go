package main

type WorkerPool struct {
	concurrency    int
	taskQuequeChan chan *Task
	runBackground  chan bool
	workers        []*Worker
}

func NewWorkerPool(concurrency int) *WorkerPool {
	return &WorkerPool{
		concurrency:    concurrency,
		taskQuequeChan: make(chan *Task, concurrency),
	}
}

func (p *WorkerPool) runWorkerPool() {
	for i := 0; i < p.concurrency; i++ {
		worker := NewWorker(p.taskQuequeChan, i)
		p.workers = append(p.workers, worker)
		go worker.Run()
	}

	p.runBackground = make(chan bool)
	<-p.runBackground
}

func (p *WorkerPool) Enqueue(task *Task) {
	p.taskQuequeChan <- task
}

func (p *WorkerPool) Stop() {
	for i := range p.workers {
		p.workers[i].Stop()
	}
	p.runBackground <- true
}
