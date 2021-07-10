package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Dmitry-Ship/playground/chat"
	"github.com/Dmitry-Ship/playground/workerpool"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	maxWorkers, err := strconv.Atoi(os.Getenv("MAX_WORKERS"))

	if err != nil {
		log.Fatal("Error reading MAX_WORKERS")
	}

	dispatcher := workerpool.NewDispatcher(maxWorkers)
	go dispatcher.Run()

	http.HandleFunc("/", chat.RootHandler)
	http.Handle("/socket", websocket.Handler(chat.SocketHandler))

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	fmt.Println("Listening to: http://" + host + ":" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
