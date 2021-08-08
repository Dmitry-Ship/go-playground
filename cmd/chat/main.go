package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dmitry-Ship/playground/pkg/chat"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", chat.RootHandler)
	http.Handle("/socket", websocket.Handler(chat.SocketHandler))

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	fmt.Println("Listening to: http://" + host + ":" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
