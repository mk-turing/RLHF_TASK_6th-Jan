package _12277

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func BenchmarkServerPerformance(b *testing.B) {
	// Get the number of workers from the environment variable
	numWorkersStr := os.Getenv("NUM_WORKERS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		b.Fatalf("Invalid NUM_WORKERS environment variable: %v", err)
	}

	// Create a worker pool with the specified number of workers
	workers := make(chan struct{}, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers <- struct{}{}
	}

	// Simulate work using a goroutine
	go func() {
		for {
			<-workers
			time.Sleep(time.Millisecond * 10) // Simulate work duration
			workers <- struct{}{}
		}
	}()

	// Run the benchmark loop
	for n := 0; n < b.N; n++ {
		// Simulate client requests by acquiring a worker
		<-workers
		// Release the worker after simulating work
		workers <- struct{}{}
	}
}
