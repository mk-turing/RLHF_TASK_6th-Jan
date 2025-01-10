package _12283

import (
	"sync"
	"testing"
)

// Function to simulate heavy load
func TestHeavyLoad(t *testing.T) {
	const numGoroutines = 100
	const numIterations = 100000
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numIterations; j++ {
				// Simulate work here
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
