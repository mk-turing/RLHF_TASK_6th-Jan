package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Define the custom error type
type TaskError struct {
	TaskID    int
	Error     error
	Severity  string // e.g., "critical", "warning", "info"
	Timestamp time.Time
}

func (te *TaskError) Error() string {
	return fmt.Sprintf("task %d failed: %v - severity: %s, timestamp: %v", te.TaskID, te.Error, te.Severity, te.Timestamp)
}

// Define the task function that may encounter errors and uses the custom error type
func taskWithCustomError(id int, errorCh chan TaskError) {
	startTime := time.Now()
	// Simulate work with a random delay
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	// Introduce a random error
	if rand.Intn(10) == 0 {
		err := fmt.Errorf("random error occurred")
		errorCh <- TaskError{TaskID: id, Error: err, Severity: "warning", Timestamp: startTime}
		return
	}
	// No error, return nil
	errorCh <- TaskError{TaskID: id, Error: nil, Severity: "info", Timestamp: startTime}
}

func taskTraditionalWithCustomError(id int) TaskError {
	startTime := time.Now()
	// Simulate work with a random delay
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	// Introduce a random error
	if rand.Intn(10) == 0 {
		err := fmt.Errorf("random error occurred")
		return TaskError{TaskID: id, Error: err, Severity: "warning", Timestamp: startTime}
	}
	// No error, return nil
	return TaskError{TaskID: id, Error: nil, Severity: "info", Timestamp: startTime}
}

type WorkerPool struct {
	minWorkers int
	maxWorkers int
	workItems  chan int
	workers    int
	done       chan struct{}
	errorCh    chan TaskError
	mutex      sync.Mutex
	wg         sync.WaitGroup
}

func NewWorkerPool(min, max int) *WorkerPool {
	return &WorkerPool{
		minWorkers: min,
		maxWorkers: max,
		workItems:  make(chan int),
		done:       make(chan struct{}),
		errorCh:    make(chan TaskError, max),
	}
}
// Rest of the code remains the same

func main() {
	// ... (Rest of the code remains the same)
	fmt.Println("\n*** Testing with Custom Error Type ***")

	// Test using channels with custom error type
	for _, workload := range workloads {
		wpChannelsCustomError := NewWorkerPool(minWorkers, maxWorkers)
		wpChannelsCustomError.Start()
		wpChannelsCustomError.Submit(workload)
		// Rest of the code remains the same

		// Test using traditional return values with custom error type
		for _, workload := range workloads {
			wpTraditionalCustomError := NewWorkerPool(minWorkers, maxWorkers)
			wpTraditionalCustomError.Start()
			wpTraditionalCustomError.Submit(workload)
			// Rest of the code remains the same
		}
	}