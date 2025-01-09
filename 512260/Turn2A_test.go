package _12260

import (
	"sync"
	"testing"
)

func add(a, b int) int {
	return a + b
}

// BenchmarkAddConcurrently demonstrates concurrent benchmarking
func BenchmarkAddConcurrently(b *testing.B) {
	var wg sync.WaitGroup
	var totalSum int

	// Define a worker function that performs the addition
	addWorker := func() {
		for i := 0; i < b.N; i++ {
			totalSum += add(1, 2)
		}
		wg.Done()
	}

	// Spawn multiple goroutines to perform the task concurrently
	const numGoroutines = 4
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go addWorker()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Measure the total time taken for all goroutines
	b.StopTimer()
	b.Logf("Total sum: %d", totalSum)
}
