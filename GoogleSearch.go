package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = FakeSearch("web")
	Image = FakeSearch("image")
	Video = FakeSearch("video")
)

type Result string

type Search func(query string) Result

func FakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func GoogleSequential(query string) (results []Result) {
	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))
	return
}

func GoogleConcurrent(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	for i := 0; i < 3; i++ {
		result := <-c
		results = append(results, result)
	}
	return
}

func GoogleNoLocks(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

var (
	Web1   = FakeSearch("web")
	Image1 = FakeSearch("image")
	Video1 = FakeSearch("video")
	Web2   = FakeSearch("web")
	Image2 = FakeSearch("image")
	Video2 = FakeSearch("video")
	Web3   = FakeSearch("web")
	Image3 = FakeSearch("image")
	Video3 = FakeSearch("video")
)

func GoogleReplicas(query string) (results []Result) {
	c := make(chan Result)
	go func() { c <- First(query, Web, Web1, Web2, Web3) }()
	go func() { c <- First(query, Image, Image1, Image2, Image3) }()
	go func() { c <- First(query, Video, Video1, Video2, Video3) }()

	timeout := time.After(100 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

func TestSearch() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	result := GoogleReplicas("golang")
	elapsed := time.Since(start)
	fmt.Println(result)
	fmt.Println(elapsed)
}
