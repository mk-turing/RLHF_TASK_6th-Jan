package _12277

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Worker represents a server process
type Worker struct {
	ID       int
	Requests chan struct{}
	wg       *sync.WaitGroup
	counter  *int64
}

// NewWorker creates a new Worker
func NewWorker(id int, wg *sync.WaitGroup, counter *int64) *Worker {
	return &Worker{
		ID:       id,
		Requests: make(chan struct{}),
		wg:       wg,
		counter:  counter,
	}
}

// Start starts the Worker's goroutine
func (w *Worker) Start() {
	go func() {
		for range w.Requests {
			// Simulate work duration
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			atomic.AddInt64(w.counter, 1)
			w.wg.Done()
		}
	}()
}

// BenchmarkServerPerformance benchmarks server performance with varying levels of parallelism
func BenchmarkServerPerformance(b *testing.B) {
	// Get the number of workers from the environment variable
	numWorkersStr := os.Getenv("NUM_WORKERS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		b.Fatalf("Invalid NUM_WORKERS environment variable: %v", err)
	}

	// Get the load distribution pattern from the environment variable
	loadPatternStr := os.Getenv("LOAD_PATTERN")
	loadPattern := 0 // Default to round-robin
	switch loadPatternStr {
	case "round-robin":
		loadPattern = 0
	case "random":
		loadPattern = 1
	default:
		b.Fatalf("Invalid LOAD_PATTERN environment variable: %s", loadPatternStr)
	}

	var wg sync.WaitGroup
	var counter int64 // Counter using atomic operations
	workers := make([]*Worker, numWorkers)

	// Initialize and start workers
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, &wg, &counter)
		workers[i] = worker
		worker.Start()
	}

	// Log the initial setup
	log.Printf("Benchmark started: %s", b.Name())
	log.Printf("Number of workers: %d", numWorkers)
	log.Printf("Load pattern: %s", loadPatternStr)

	// Simulate different traffic patterns based on the load pattern
	startTime := time.Now()
	for n := 0; n < b.N; n++ {
		switch loadPattern {
		case 0: // Round-robin
			worker := workers[n%numWorkers]
			worker.Requests <- struct{}{}
			wg.Add(1) // Add work before sending the request
		case 1: // Random
			worker := workers[rand.Intn(numWorkers)]
			worker.Requests <- struct{}{}
			wg.Add(1) // Add work before sending the request
		}
	}

	// Wait for all workers to complete their requests
	wg.Wait()

	// Calculate the duration of the benchmark
	duration := time.Since(startTime)

	// Log performance metrics
	log.Printf("Benchmark completed: %s", b.Name())
	log.Printf("Benchmark duration: %v", duration)
	log.Printf("Total requests processed: %d", atomic.LoadInt64(&counter))

	if duration > 0 {
		opsPerSecond := float64(b.N) / duration.Seconds()
		log.Printf("Operations per second: %.2f", opsPerSecond)
	} else {
		log.Println("Benchmark duration is 0, unable to calculate operations per second.")
	}
}

func Example() {
	// Set the number of workers to 4 and use random load distribution
	os.Setenv("NUM_WORKERS", "4")
	os.Setenv("LOAD_PATTERN", "random")

	// Run the benchmark with a duration of 5 seconds
	testing.Benchmark(BenchmarkServerPerformance)
}
