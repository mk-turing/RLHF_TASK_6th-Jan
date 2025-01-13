package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type logMessage func(level, message string) error

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
			err := processor("Info", fmt.Sprintf("Goroutine %d logged this informational message.", id))
			if err != nil {
				fmt.Printf("Error logging for goroutine %d: %v\n", id, err)
			}
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
			err := processor("Info", fmt.Sprintf("Goroutine %d modified log level and logged this informational message.", id))
			if err != nil {
				fmt.Printf("Error logging for goroutine %d: %v\n", id, err)
			}
		}(i + 5)
	}

	wg.Wait()
}

// Create a log processor function using a closure
func createLogProcessor() logMessage {
	var logCounter int = 0
	return func(level, message string) error {
		mutex.Lock()
		defer mutex.Unlock()

		if err := logInternal(level, message); err != nil {
			return err
		}

		logCounter++
		if logCounter <= levelLimit {
			fmt.Printf("%v: %v\n", time.Now(), message)
		} else {
			return fmt.Errorf("Log limit reached for message: %s", message)
		}

		return nil
	}
}

// Simulated internal logging function that can fail
func logInternal(level, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(200*time.Millisecond))
	defer cancel()

	select {
	case <-ctx.Done():
		return fmt.Errorf("Internal logging timed out")
	default:
		// Simulate some logging logic
		time.Sleep(time.Duration(100 * time.Millisecond))
		return nil
	}
}
