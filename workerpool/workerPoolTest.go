package workerpool

import (
	"fmt"
	"sync"

	"github.com/Dmitry-Ship/playground/common"
)

func mockWork(dispatcher *Dispatcher) {
	c := make(chan Result)
	var wg sync.WaitGroup
	for i := 0; i < 40; i++ {
		job := NewJob(i, c)
		wg.Add(1)

		go dispatcher.Enqueue(job)

	}

	go func() {
		for r := range c {
			fmt.Printf("✅ result: %d \n", r.JobId)
			wg.Done()
		}
	}()

	wg.Wait()
}

func program() {
	const concurrency = 10
	dispatcher := NewDispatcher(concurrency)
	go dispatcher.Run()
	mockWork(dispatcher)
}

func TestWorkerPool() {
	common.MeasureTime(program)
}
