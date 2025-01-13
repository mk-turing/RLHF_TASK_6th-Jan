package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	count int
	mutex sync.Mutex
}

// Increment method for Counter
func (c *Counter) Increment() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.count++
}

// Value method for Counter to retrieve the current value
func (c *Counter) Value() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.count
}

// Callback function type definition
type Callback func()

// Increment counter function to pass as callback
func incrementCounter(id int, counter *Counter, callback Callback) {
	defer wg.Done()

	// Calling the callback (Increment method)
	callback()

	// Print the incremented value (to simulate real-world task completion)
	fmt.Printf("Incremented by goroutine %d: %d\n", id, counter.Value())
}

var wg sync.WaitGroup

func main() {
	// Create an instance of Counter
	counter := &Counter{}

	// Function pointer as callback (Counter.Increment)
	callback := counter.Increment

	// Start multiple goroutines (10 in this case)
	wg.Add(10)
	for i := 1; i <= 10; i++ {
		go incrementCounter(i, counter, callback)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Print final counter value
	fmt.Println("Final counter value:", counter.Value())
}
