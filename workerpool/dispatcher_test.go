package workerpool

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func performLongWork() (string, error) {
	randomNumber := rand.Intn(100)
	return strconv.Itoa(randomNumber), nil
}

func createJobs(number int) []Job {
	jobs := []Job{}
	c := make(chan Result)

	for n := 0; n < number; n++ {
		job := NewJob(n, c, performLongWork)
		jobs = append(jobs, job)
	}
	return jobs
}

func BenchmarkConcurrent(b *testing.B) {
	dispatcher := NewDispatcher(100) // start up worker pool
	go dispatcher.Run()

	for n := 0; n < b.N; n++ {
		c := make(chan Result)
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			job := NewJob(i, c, performLongWork)
			wg.Add(1)

			go dispatcher.Enqueue(job)

		}

		go func() {
			for range c {
				// fmt.Printf("✅ result: %d \n", r.JobId)
				// <-r
				wg.Done()
			}
		}()

		wg.Wait()
	}
}
