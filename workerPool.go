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
	executeTask func() (error, string)
}

func (t *Task) Run(wg *sync.WaitGroup) string {
	err, result := t.executeTask()
	t.Err = err
	// wg.Done()

	return result
}

type Result struct {
	Value  string
	TaskId int
}

type WorkerPool struct {
	limit       int
	taskQueue   []*Task
	tasksChan   chan *Task
	resultsChan chan *Result
	wg          sync.WaitGroup
}

func NewWorkerPool(numberOfWorkers int, capacity int) *WorkerPool {
	return &WorkerPool{
		limit:       numberOfWorkers,
		tasksChan:   make(chan *Task, capacity),
		resultsChan: make(chan *Result, capacity),
	}
}

func (p *WorkerPool) addTask(task *Task) {
	p.taskQueue = append(p.taskQueue, task)
}

func (p *WorkerPool) runWorker() {
	p.wg.Add(1)

	for task := range p.tasksChan {

		result := Result{
			Value:  task.Run(&p.wg),
			TaskId: task.Id,
		}

		p.resultsChan <- &result
	}

	p.wg.Done()
}

func (p *WorkerPool) run() {
	for i := 0; i < p.limit; i++ {
		go p.runWorker()
	}

	go func() {
		for result := range p.resultsChan {
			fmt.Println(fmt.Sprintf("task id: %d, value: %s", result.TaskId, result.Value))
		}
	}()

	for _, task := range p.taskQueue {
		p.tasksChan <- task
	}

	close(p.tasksChan)
	p.wg.Wait()
	close(p.resultsChan)
}

func performLongWork() (error, string) {
	randomNumber := rand.Intn(10)

	time.Sleep(time.Duration(randomNumber) * time.Second)
	return nil, strconv.Itoa(randomNumber)
}

func testWorkerPool() {
	workerPool := NewWorkerPool(10, 10)

	for i := 0; i < 20; i++ {
		task := &Task{executeTask: performLongWork, Id: i}
		workerPool.addTask(task)
	}

	measureTime(workerPool.run)
}
