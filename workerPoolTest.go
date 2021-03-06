package main

import (
	"fmt"
	"math/rand"
	"time"
)

func readResults(results chan *Result) {
	for result := range results {
		fmt.Printf("✅ finished %d\n", result.TaskId)

	}
}

func mockWork(workerPool *WorkerPool, results chan *Result) {
	for {
		taskID := rand.Intn(100)

		if taskID == 0 {
			workerPool.Stop()
		}

		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
		task := NewTask(taskID, results)
		workerPool.Enqueue(task)
	}

}

func program() {
	const concurrency = 10
	results := make(chan *Result, concurrency)
	workerPool := NewWorkerPool(concurrency)

	go mockWork(workerPool, results)
	go readResults(results)

	workerPool.runWorkerPool()
}

func testWorkerPool() {
	measureTime(program)
}
