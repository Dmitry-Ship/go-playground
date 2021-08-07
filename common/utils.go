package common

import (
	"fmt"
	"net/http"
	"os"
	"sync"
)

func ErrorHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if _, ok := recover().(os.LinkError); ok {
				w.WriteHeader(500)
			}
		}()
		fn(w, r)
	}
}

func FanIn(inputChannels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	outputChannel := make(chan string)

	wg.Add(len(inputChannels))
	for _, inputChannel := range inputChannels {
		go func(ic <-chan string) {
			for message := range ic {
				outputChannel <- message
			}
			wg.Done()
		}(inputChannel)
	}

	go func() {
		wg.Wait()
		close(outputChannel)
		fmt.Println("done: ")
	}()
	return outputChannel
}
