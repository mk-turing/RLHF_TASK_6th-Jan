package main

import (
	"fmt"
	"sync"
	"time"
)

type logMessage func(level, message string)

var totalLogs int = 0
var levelLimit int = 5

var wg sync.WaitGroup
var mutex sync.Mutex

func main() {
	// Create a log processor using a closure
	processor := createLogProcessor()

	wg.Add(10)
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			processor("Info", fmt.Sprintf("Goroutine %d logged this informational message.", id))
			time.Sleep(time.Duration(500 * time.Millisecond))
		}(i)
	}

	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			if id%2 == 0 {
				mutex.Lock()
				levelLimit--
				mutex.Unlock()
			}
			processor("Info", fmt.Sprintf("Goroutine %d modified log level and logged this informational message.", id))
		}(i + 5)
	}

	wg.Wait()
}

// Create a log processor function using a closure
func createLogProcessor() logMessage {
	var logCounter int = 0
	return func(level, message string) {
		mutex.Lock()
		defer mutex.Unlock()

		logCounter++
		if logCounter <= levelLimit {
			fmt.Printf("%v: %v\n", time.Now(), message)
		}
	}
}
