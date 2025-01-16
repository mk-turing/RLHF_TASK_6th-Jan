
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Maximum number of concurrent tasks
	maxWorkers = 100
)

// Define the task function that may encounter errors
func task(id int, wg *sync.WaitGroup, errorCh chan error) {
	defer wg.Done()
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

// Define the task function that may encounter errors using traditional return values
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

func worker(id int, workCh <-chan int, errorCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	for taskID := range workCh {
		// Simulate work with a random delay
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		// Introduce a random error
		if rand.Intn(10) == 0 {
			errorCh <- fmt.Errorf("task %d failed", taskID)
			continue
		}
		// No error, continue with the next task
	}
}

func adjustWorkers(numTasks int, numWorkers int32, workCh chan int, errorCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		// Calculate the desired number of workers based on the number of tasks
		desiredWorkers := numTasks / 100
		if desiredWorkers > maxWorkers {
			desiredWorkers = maxWorkers
		}

		// Get the current number of workers
		currentWorkers := atomic.LoadInt32(&numWorkers)

		if desiredWorkers > currentWorkers {
			// Start more workers if needed
			for i := currentWorkers; i < desiredWorkers; i++ {
				go worker(int(i), workCh, errorCh, wg)
				atomic.AddInt32(&numWorkers, 1)
			}
		} else if desiredWorkers < currentWorkers {
			// Close the work channel to signal workers to stop
			close(workCh)
			wg.Wait()

			// Adjust the number of workers and start again
			atomic.StoreInt32(&numWorkers, int32(desiredWorkers))
			workCh = make(chan int, numTasks)
			wg.Add(desiredWorkers)
			for i := 0; i < desiredWorkers; i++ {
				go worker(int(i), workCh, errorCh, wg)
			}
		}

		// Wait for a little while before adjusting again
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	var numTasks int
	fmt.Print("Enter the number of tasks: ")
	fmt.Scan(&numTasks)

	// Using channels for error reporting
	errorCh := make(chan error, numTasks)
	var wg sync.WaitGroup
	var numWorkers int32
	atomic.StoreInt32(&numWorkers, 0)