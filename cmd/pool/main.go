package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/Dmitry-Ship/playground/pkg/workerpool"
	"github.com/joho/godotenv"
)

func performLongWork() (string, error) {
	randomNumber := rand.Intn(10000)
	time.Sleep(time.Duration(randomNumber) * time.Millisecond)
	return strconv.Itoa(randomNumber), nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))

	if err != nil {
		log.Fatal("Error reading MAX_WORKERS")
	}

	dispatcher := workerpool.NewDispatcher(maxWorkers) // start up worker pool
	go dispatcher.Run()

	resultChan := make(chan workerpool.Result)

	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)
	id := 0
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:

				for i := 0; i < rand.Intn(50); i++ {
					id++
					job := workerpool.NewJob(id, resultChan, performLongWork)
					time.Sleep(time.Duration(100) * time.Millisecond)
					go dispatcher.Enqueue(job)

				}
			}
		}
	}()

	for result := range resultChan {
		fmt.Println("✅finished job", result.JobId)
	}
}
