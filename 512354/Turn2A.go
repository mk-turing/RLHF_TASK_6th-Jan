package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	counter = 0
	mu      = &sync.Mutex{}
)

func asyncOperation(callback func() error, delay time.Duration) {
	time.Sleep(delay)
	if err := callback(); err != nil {
		fmt.Println("Async operation failed:", err)
	}
}

func incrementCounter() error {
	mu.Lock()
	defer mu.Unlock()

	counter++

	// Simulate an error condition
	if counter%3 == 0 {
		return fmt.Errorf("counter reached multiple of 3: %d", counter)
	}
	return nil
}

func main() {
	for i := 0; i < 10; i++ {
		go asyncOperation(incrementCounter, time.Duration(i)*time.Second)
	}
	time.Sleep(11 * time.Second)
	fmt.Println("Final counter value:", counter)
}
