package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Define the task function that may encounter errors
func task(id int, errorCh chan error) {
	// Simulate work with a random delay
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	// Introduce a random error
	if rand.Intn(10) == 0 {
		errorCh <- fmt.Errorf("task %d failed", id)
		return
	}
	// No error, return nil
	errorCh <- nil
}

func taskTraditional(id int) error {
	// Simulate work with a random delay
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	// Introduce a random error
	if rand.Intn(10) == 0 {
		return fmt.Errorf("task %d failed", id)
	}
	// No error, return nil
	return nil
}

// Dynamic worker pool
type WorkerPool struct {
	minWorkers int
	maxWorkers int
	workItems  chan int
	workers    int
	done       chan struct{}
	errorCh    chan error
	mutex      sync.Mutex
	wg         sync.WaitGroup
}

func NewWorkerPool(min, max int) *WorkerPool {
	return &WorkerPool{
		minWorkers: min,
		maxWorkers: max,
		workItems:  make(chan int),
		done:       make(chan struct{}),
		errorCh:    make(chan error, max),
	}
}

func (wp *WorkerPool) Start() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	wp.workers = wp.minWorkers
	for i := 0; i < wp.workers; i++ {
		go wp.worker()
	}
}

func (wp *WorkerPool) Stop() {
	wp.wg.Wait() // Wait for all workers to finish
	close(wp.done)
	close(wp.workItems) // Close the work items channel
	close(wp.errorCh)   // Close the error channel
}

func (wp *WorkerPool) worker() {
	for {
		select {
		case <-wp.done:
			return
		case item, ok := <-wp.workItems:
			if !ok {
				return
			}
			wp.wg.Add(1)           // Ensure Add() is called for each task in the worker
			task(item, wp.errorCh) // Run task with error channel
			wp.wg.Done()           // Ensure Done() is called for each task after it finishes
		}
	}
}

func (wp *WorkerPool) Submit(numTasks int) {
	for i := 0; i < numTasks; i++ {
		wp.workItems <- i
	}
}

func (wp *WorkerPool) GetErrors() []error {
	var errors []error
	for err := range wp.errorCh {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func main() {
	rand.Seed(time.Now().UnixNano())
	minWorkers := 1
	maxWorkers := 100
	workloads := []int{10, 25, 50, 75, 100}

	// Test using channels
	for _, workload := range workloads {
		wpChannels := NewWorkerPool(minWorkers, maxWorkers)
		wpChannels.Start()
		wpChannels.Submit(workload)

		startChannels := time.Now()
		wpChannels.Stop()
		endChannels := time.Now()

		timeTakenChannels := endChannels.Sub(startChannels)
		errorsChannels := wpChannels.GetErrors()

		fmt.Printf("Workload: %d, Time taken using channels: %v, Errors: %d\n", workload, timeTakenChannels, len(errorsChannels))
	}

	// Test using traditional return values
	for _, workload := range workloads {
		wpTraditional := NewWorkerPool(minWorkers, maxWorkers)
		wpTraditional.Start()
		wpTraditional.Submit(workload)

		startTraditional := time.Now()
		wpTraditional.Stop()
		endTraditional := time.Now()

		timeTakenTraditional := endTraditional.Sub(startTraditional)
		errorsTraditional := wpTraditional.GetErrors()

		fmt.Printf("Workload: %d, Time taken using traditional return values: %v, Errors: %d\n", workload, timeTakenTraditional, len(errorsTraditional))
	}
}
