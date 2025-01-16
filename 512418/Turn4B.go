package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Define a custom error type with additional context
type CustomError struct {
	TaskID    int
	ErrorMsg  string
	Severity  string // e.g., "info", "warning", "error"
	TimeTaken time.Duration
}

func (ce CustomError) Error() string {
	return fmt.Sprintf("Task ID: %d, Error: %s, Severity: %s, Time Taken: %s", ce.TaskID, ce.ErrorMsg, ce.Severity, ce.TimeTaken.String())
}

// Global error handler with buffered channel
var globalErrorHandler = newErrorHandler(100)

type errorHandler struct {
	errorCh chan *CustomError
}

func newErrorHandler(buffSize int) *errorHandler {
	return &errorHandler{
		errorCh: make(chan *CustomError, buffSize),
	}
}

func (eh *errorHandler) LogAndRetry(err *CustomError) {
	// Simulate logging or retry logic
	fmt.Println("Logged Error:", err)
	// retry logic could be implemented here using another goroutine
}

func (eh *errorHandler) Start() {
	go func() {
		for err := range eh.errorCh {
			eh.LogAndRetry(err)
		}
	}()
}

func task(id int, taskCh chan int) {
	defer func() {
		if r := recover(); r != nil {
			err := &CustomError{
				TaskID:   id,
				ErrorMsg: fmt.Sprintf("%v", r),
				Severity: "panic",
			}
			globalErrorHandler.errorCh <- err
		}
	}()
	// Simulate work with a random delay
	//startTime := time.Now()
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	//timeTaken := time.Since(startTime)

	// Introduce a random error
	if rand.Intn(10) == 0 {
		panic("task failed")
	}

	taskCh <- id // Signal task completion
}

func main() {
	rand.Seed(time.Now().UnixNano())
	globalErrorHandler.Start()

	numTasks := 1000
	taskCh := make(chan int, numTasks)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < 10; i++ { // Adjust the number of workers as needed
		wg.Add(1)
		go func() {
			defer wg.Done()
			for taskID := range taskCh {
				task(taskID, taskCh)
			}
		}()
	}

	// Submit tasks
	for i := 0; i < numTasks; i++ {
		taskCh <- i
	}
	close(taskCh)

	// Wait for all tasks to complete
	wg.Wait()

	fmt.Println("All tasks completed.")
}
