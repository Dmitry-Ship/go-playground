package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Dmitry-Ship/playground/chat"
	"github.com/Dmitry-Ship/playground/workerpool"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func performLongWork() (string, error) {
	randomNumber := rand.Intn(1000)
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

	// var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		job := workerpool.NewJob(i, resultChan, performLongWork)
		// wg.Add(1)

		go dispatcher.Enqueue(job)

	}

	go func() {
		for result := range resultChan {
			fmt.Println(result.JobId)
			// wg.Done()
		}
	}()

	http.HandleFunc("/", chat.RootHandler)
	http.Handle("/socket", websocket.Handler(chat.SocketHandler))

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	fmt.Println("Listening to: http://" + host + ":" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
