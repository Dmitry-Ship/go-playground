package main

import (
	"fmt"
	"net/http"
	"time"
)

type SearchResult struct {
	value string
}

func Web(query string) SearchResult {
	_, err := http.Get("https://www.google.com/search?q=" + query)
	if err != nil {

		return SearchResult{value: "result"}
	}

	return SearchResult{value: ""}
}

func Image(query string) SearchResult {
	_, err := http.Get("https://www.google.com/search?q=" + query)
	if err != nil {

		return SearchResult{value: "result"}
	}

	return SearchResult{value: ""}
}

func Video(query string) SearchResult {
	_, err := http.Get("https://www.google.com/search?q=" + query)
	if err != nil {
		return SearchResult{value: "result"}
	}

	return SearchResult{value: ""}
}

func Google(query string) []SearchResult {
	var results = []SearchResult{}

	c := make(chan SearchResult)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	for i := 0; i < 1; i++ {
		results = append(results, <-c)
	}

	return results
}

func testGoogle() {
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
