package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
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

func main() {
	numTasks := 1000
	// Using channels for error reporting
	errorCh := make(chan error, numTasks)
	var wg sync.WaitGroup
	wg.Add(numTasks)

	// Run tasks concurrently
	for i := 0; i < numTasks; i++ {
		go task(i, &wg, errorCh)
	}

	// Measure the time taken using channels
	start := time.Now()
	wg.Wait()
	close(errorCh)
	end := time.Now()
	timeTakenChannels := end.Sub(start)

	// Using traditional return values
	var errors []error
	var wg2 sync.WaitGroup
	wg2.Add(numTasks)

	for i := 0; i < numTasks; i++ {
		go func(id int) {
			defer wg2.Done()
			err := taskTraditional(id)
			if err != nil {
				errors = append(errors, err)
			}
		}(i)
	}

	// Measure the time taken using traditional return values
	start2 := time.Now()
	wg2.Wait()
	end2 := time.Now()
	timeTakenTraditional := end2.Sub(start2)

	// Print the results
	fmt.Printf("Number of tasks: %d\n", numTasks)
	fmt.Printf("Time taken using channels: %v\n", timeTakenChannels)
	fmt.Printf("Time taken using traditional return values: %v\n", timeTakenTraditional)

	// Check if any errors occurred using channels
	for err := range errorCh {
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Check if any errors occurred using traditional return values
	if len(errors) > 0 {
		fmt.Println("Errors occurred using traditional return values:")
		for _, err := range errors {
			fmt.Printf("Error: %v\n", err)
		}
	}
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
