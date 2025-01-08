package main

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// A shared resource to simulate a critical section
var sharedResource = &SharedResource{
	value: 0,
	mu:    &sync.Mutex{}, // Initialize the Mutex
}

// A simple shared resource struct
type SharedResource struct {
	value int64
	mu    *sync.Mutex
}

func (s *SharedResource) increment() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value++
	return s.value
}

// A function to be benchmarked with concurrency
func concurrentFunction() {
	sharedResource.increment()
}

// A helper function to run multiple goroutines concurrently
func runConcurrently(numGoroutines int, duration time.Duration, f func()) {
	numWorkers := numGoroutines
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	timer := time.After(duration)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			select {
			case <-timer:
				return
			default:
				f()
			}
		}()
	}

	wg.Wait()
}

// BenchmarkConcurrentFunction measures the performance of the function under concurrent load
func BenchmarkConcurrentFunction(b *testing.B) {
	const numGoroutines = 100         // Adjust the number of goroutines as needed
	const duration = time.Second * 10 // Test duration

	// Initialize counts with a length of 5 (as you are running the test 5 times)
	counts := make([]int64, 5)

	for i := 0; i < 5; i++ { // Repeat the test to ensure reproducibility
		sharedResource.value = 0 // Reset shared resource for each test run
		atomic.StoreInt64(&counts[i], 0)

		go func(idx int) {
			for i := 0; i < b.N; i++ {
				counts[idx] += sharedResource.increment()
			}
		}(i)

		runConcurrently(numGoroutines, duration, concurrentFunction)

		total := atomic.LoadInt64(&counts[i])

		// Report the metric for ops/sec
		opsPerSec := float64(total) / duration.Seconds()
		b.ReportMetric(opsPerSec, "ops/sec")
	}
}
