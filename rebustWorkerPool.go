package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct {
	Err         error
	Id          int
	executeTask func() error
}

func (t *Task) Run(wg *sync.WaitGroup) {
	t.Err = t.executeTask()
	fmt.Println(fmt.Sprintf("task id %d, executed", t.Id))
	wg.Done()
}

type WorkerPool struct {
	numberOfWorkers int
	tasks           []*Task
	tasksChan       chan *Task
	wg              sync.WaitGroup
}

func NewWorkerPool(numberOfWorkers int) *WorkerPool {
	return &WorkerPool{
		numberOfWorkers: numberOfWorkers,
		tasksChan:       make(chan *Task),
	}
}

func (p *WorkerPool) addTask(task *Task) {
	p.tasks = append(p.tasks, task)
}

func (p *WorkerPool) run() {
	p.wg.Add(len(p.tasks))

	for i := 0; i < p.numberOfWorkers; i++ {
		go p.runWorker()
	}

	for _, task := range p.tasks {
		p.tasksChan <- task
	}

	// all workers return
	close(p.tasksChan)

	p.wg.Wait()
}

func (p *WorkerPool) runWorker() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
	}
}

func performLongWork() error {
	randomNumber := rand.Intn(10)

	time.Sleep(time.Duration(randomNumber) * time.Second)
	return nil
}

func testRobustWorkerPool() {
	workerPool := NewWorkerPool(20)

	for i := 0; i < 100; i++ {
		task := &Task{executeTask: performLongWork, Id: i}
		workerPool.addTask(task)
	}

	measureTime(workerPool.run)
}
