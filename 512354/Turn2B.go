package main

import (
	"fmt"
	"sync"
	"time"
)

type OperationError struct {
	Err error
}

func (oe *OperationError) Error() string {
	return fmt.Sprintf("Operation failed: %v", oe.Err)
}

var counter = 0
var mu sync.Mutex

func asyncOperationWithError(id int, callback func(int, error), delay time.Duration, errChannel chan error) {
	time.Sleep(delay)

	// Simulate an error in 50% of the cases
	if id%2 == 0 {
		err := fmt.Errorf("simulated error for operation %d", id)
		errChannel <- &OperationError{Err: err}
		return
	}

	mu.Lock()
	counter++
	mu.Unlock()

	callback(id, nil)
}

func main() {
	done := make(chan struct{})
	errChannel := make(chan error)

	for i := 0; i < 10; i++ {
		go asyncOperationWithError(i, func(id int, err error) {
			if err != nil {
				fmt.Printf("Callback %d received error: %v\n", id, err)
				return
			}
			fmt.Printf("Callback %d successfully completed.\n", id)
		}, time.Duration(i)*time.Second, errChannel)
	}

	go func() {
		for err := range errChannel {
			fmt.Printf("Error Received: %v\n", err)
		}
		close(done)
	}()

	time.Sleep(11 * time.Second)

	fmt.Println("Final counter value:", counter)

	<-done
	fmt.Println("All operations complete.")
}
