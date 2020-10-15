package main

import (
	"fmt"
	"math/rand"
	"strconv"
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
	wg.Done()
}

type Result struct {
	Value string
}

type WorkerPool struct {
	numberOfWorkers int
	tasks           []*Task
	tasksChan       chan *Task
	resultsChan     chan *Result
	wg              sync.WaitGroup
}

func NewWorkerPool(numberOfWorkers int) *WorkerPool {
	return &WorkerPool{
		numberOfWorkers: numberOfWorkers,
		tasksChan:       make(chan *Task),
		resultsChan:     make(chan *Result),
	}
}

func (p *WorkerPool) addTask(task *Task) {
	p.tasks = append(p.tasks, task)
}

func (p *WorkerPool) runWorker() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
		result := Result{Value: "task id " + strconv.Itoa(task.Id)}

		p.resultsChan <- &result
	}
}

func (p *WorkerPool) run() {
	p.wg.Add(len(p.tasks))

	for i := 0; i < p.numberOfWorkers; i++ {
		go p.runWorker()
	}

	go func() {
		for _, task := range p.tasks {
			p.tasksChan <- task
		}
	}()

	go func() {
		for result := range p.resultsChan {
			fmt.Println("result: " + result.Value)
		}
	}()

	p.wg.Wait()
	close(p.tasksChan)
	close(p.resultsChan)
}

func performLongWork() error {
	randomNumber := rand.Intn(20)

	time.Sleep(time.Duration(randomNumber) * time.Second)
	return nil
}

func testWorkerPool() {
	workerPool := NewWorkerPool(20)

	for i := 0; i < 100; i++ {
		task := &Task{executeTask: performLongWork, Id: i}
		workerPool.addTask(task)
	}

	measureTime(workerPool.run)
}
