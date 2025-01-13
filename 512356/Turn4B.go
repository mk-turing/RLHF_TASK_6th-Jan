package main

import (
	"errors"
	"fmt"
	"io"
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
				fmt.Printf("Error processing log: %v\n", err)
			}
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
				fmt.Printf("Error processing log: %v\n", err)
			}
		}(i + 5)
	}

	wg.Wait()
}

// Create a log processor function using a closure
func createLogProcessor() logMessage {
	var logCounter int = 0
	var logFile *io.Writer

	// Initialize the log file (simulating resource initialization)
	logFile, err := createLogFile("application.log")
	if err != nil {
		return func(level, message string) error {
			return err
		}
	}

	// Return the closure that captures logCounter and logFile
	return func(level, message string) error {
		mutex.Lock()
		defer mutex.Unlock()

		logCounter++
		if logCounter <= levelLimit {
			_, err := fmt.Fprintf(*logFile, "%v: %v\n", time.Now(), message)
			if err != nil {
				return err // Return the error to the caller
			}
		}

		return nil // Return no error
	}
}

// Function to simulate creating a log file
func createLogFile(filename string) (*io.Writer, error) {
	// Simulate an error during file creation
	if filename == "application.log" {
		return nil, errors.New("error creating log file")
	}
	// Return a mock io.Writer
	return &io.Writer{}, nil
}
