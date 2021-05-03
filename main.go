package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Dmitry-Ship/playground/chat"
	"github.com/Dmitry-Ship/playground/workerpool"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func boring(msg string) <-chan string { // Returns receive-only channel of strings.
	c := make(chan string)
	go func() { // We launch the goroutine from inside the function.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Return the channel to the caller.
}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	workerpool.TestWorkerPool()

	// maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))

	// if err == nil {
	// 	workerpool.TestWorkerPool()

	// 	const concurrency = 5
	// 	dispatcher := workerpool.NewDispatcher(maxWorkers)
	// 	go dispatcher.Run()
	// }

	http.HandleFunc("/", chat.RootHandler)
	http.Handle("/socket", websocket.Handler(chat.SocketHandler))
	http.HandleFunc("/api", IndexHandler)
	http.HandleFunc("/api/message", GetMessageHandler)

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	fmt.Println("Listening to: http://" + host + ":" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
