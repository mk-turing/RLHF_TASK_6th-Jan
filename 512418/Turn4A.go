package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// CustomError represents an enriched error structure
type CustomError struct {
	TaskID    int
	ErrorMsg  string
	Severity  string
	TimeTaken time.Duration
}

func (ce CustomError) Error() string {
	return fmt.Sprintf("Task ID: %d, Error: %s, Severity: %s, Time Taken: %s", ce.TaskID, ce.ErrorMsg, ce.Severity, ce.TimeTaken.String())
}

// Global error handler
func globalErrorHandler(errorCh chan *CustomError, wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range errorCh {
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

// Task simulates a worker task that may encounter errors
func task(id int, errorCh chan *CustomError) {
	startTime := time.Now()
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	timeTaken := time.Since(startTime)

	if rand.Intn(10) == 0 { // Simulate a random failure
		errorCh <- &CustomError{
			TaskID:    id,
			ErrorMsg:  "Task failed",
			Severity:  "error",
			TimeTaken: timeTaken,
		}
		return
	}
	// No error
	errorCh <- nil
}

// WorkerPool manages dynamic workers and task execution
type WorkerPool struct {
	workItems chan int
	errorCh   chan *CustomError
	wg        sync.WaitGroup
}

// NewWorkerPool initializes a worker pool
func NewWorkerPool(bufferSize int) *WorkerPool {
	return &WorkerPool{
		workItems: make(chan int, bufferSize),
		errorCh:   make(chan *CustomError, bufferSize),
	}
}

// Start begins worker execution
func (wp *WorkerPool) Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go wp.worker()
	}
}

// Stop waits for all workers to finish and closes the error channel
func (wp *WorkerPool) Stop() {
	wp.wg.Wait()      // Ensure all workers have finished processing
	close(wp.errorCh) // Close the error channel
}

// Submit adds tasks to the work queue
func (wp *WorkerPool) Submit(taskCount int) {
	for i := 0; i < taskCount; i++ {
		wp.workItems <- i
	}
	close(wp.workItems) // Signal workers no more tasks will be added
}

// worker processes tasks and sends errors to the error channel
func (wp *WorkerPool) worker() {
	for taskID := range wp.workItems {
		wp.wg.Add(1) // Add before processing the task
		task(taskID, wp.errorCh)
		wp.wg.Done() // Mark task as done
	}
}

// Benchmark and comparison function
func benchmark(workerCount, taskCount int) {
	wp := NewWorkerPool(taskCount)
	wp.Start(workerCount)

	// Start error handler
	errorHandlerWG := sync.WaitGroup{}
	errorHandlerWG.Add(1)
	go globalErrorHandler(wp.errorCh, &errorHandlerWG)

	// Start benchmark
	start := time.Now()
	wp.Submit(taskCount)
	wp.Stop()
	errorHandlerWG.Wait() // Wait for the error handler to finish
	end := time.Now()

	fmt.Printf("Workers: %d, Tasks: %d, Time Taken: %v\n", workerCount, taskCount, end.Sub(start))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// High-load scenarios with varying worker and task counts
	configurations := []struct {
		workers   int
		taskCount int
	}{
		{workers: 5, taskCount: 50},
		{workers: 10, taskCount: 100},
		{workers: 20, taskCount: 200},
	}

	for _, config := range configurations {
		benchmark(config.workers, config.taskCount)
	}
}
