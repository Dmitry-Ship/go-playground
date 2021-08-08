package common

import (
	"fmt"
	"sync"
)

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
