package workerpool

import "fmt"

type Worker struct {
	ID       int
	taskChan <-chan *Task

	quit chan bool
}

func NewWorker(channel <-chan *Task, ID int) *Worker {
	return &Worker{
		ID:       ID,
		taskChan: channel,
		quit:     make(chan bool),
	}
}

func (wr *Worker) Run() {
	// for task := range wr.taskChan {
	// 	fmt.Printf("🎬 Starting task %d\n", task.Id)
	// 	task.Run()
	// }
	for {
		select {
		case task := <-wr.taskChan:
			fmt.Printf("🎬 Starting task %d\n", task.Id)
			task.Run()
		case <-wr.quit:
			return
		}
	}
}

func (wr *Worker) Stop() {
	go func() {
		wr.quit <- true
	}()
}
