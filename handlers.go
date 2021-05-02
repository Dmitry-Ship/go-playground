package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

type Request struct {
	Message string
}
type response struct {
	Data string `json:"data"`
}

func processRequest(r Request) response {
	if r.Message == "" {
		return response{
			Data: "empty message",
		}
	}

	return response{
		Data: r.Message,
	}
}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	var request Request

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := processRequest(request)

	json.NewEncoder(w).Encode(&data)
}
