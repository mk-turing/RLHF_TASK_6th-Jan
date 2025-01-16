package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup
var errCount = 0
var errorCh = make(chan error)

func simulateTask() error {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	if rand.Intn(10) == 0 {
		return fmt.Errorf("simulated network error")
	}
	return nil
}

func simulateTaskWithChannel() {
	defer wg.Done()
	err := simulateTask()
	if err != nil {
		errorCh <- err
		errCount++
	}
}

func benchmarkReturnValues(numTasks int) {
	start := time.Now()
	wg.Add(numTasks)

	for i := 0; i < numTasks; i++ {
		go func() {
			defer wg.Done()
			err := simulateTask()
			if err != nil {
				errCount++
			}
		}()
	}

	wg.Wait()
	end := time.Now()
	fmt.Printf("Return Values Method: Time taken = %v, Errors = %d\n", end.Sub(start), errCount)
}

func benchmarkChannels(numTasks int) {
	start := time.Now()
	wg.Add(numTasks)

	for i := 0; i < numTasks; i++ {
		go simulateTaskWithChannel()
	}

	for i := 0; i < numTasks; i++ {
		select {
		case <-errorCh:
			errCount++
		default:
		}
	}

	wg.Wait()
	end := time.Now()
	fmt.Printf("Channels Method: Time taken = %v, Errors = %d\n", end.Sub(start), errCount)
}

func main() {
	numTasks := 100000
	rand.Seed(time.Now().UnixNano())

	benchmarkReturnValues(numTasks)
	errCount = 0
	benchmarkChannels(numTasks)
}
