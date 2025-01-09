package _12277

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Worker struct {
	ID       int
	Requests chan struct{}
	wg       *sync.WaitGroup
}

func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:       id,
		Requests: make(chan struct{}, 100), // Buffered channel for improved concurrency
		wg:       wg,
	}
}

func (w *Worker) Start() {
	go func() {
		for req := range w.Requests {
			// Simulate work duration
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			w.wg.Done()
		}
	}()
}

func BenchmarkServerPerformance(b *testing.B) {
	// Parse environment variables
	numWorkersStr := os.Getenv("NUM_WORKERS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		b.Fatalf("Invalid NUM_WORKERS environment variable: %v", err)
	}

	maxParallelismStr := os.Getenv("MAX_PARALLELISM")
	maxParallelism, err := strconv.Atoi(maxParallelismStr)
	if err != nil {
		maxParallelism = runtime.NumCPU() // Default to the number of available CPU cores
	}

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

	// Create workers
	var wg sync.WaitGroup
	workers := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, &wg)
		workers[i] = worker
		worker.Start()
	}

	// Set the maximum parallelism
	runtime.GOMAXPROCS(maxParallelism)

	// Log setup
	log.Printf("Benchmark started: %s", b.Name())
	log.Printf("Number of workers: %d", numWorkers)
	log.Printf("Max parallelism: %d", maxParallelism)
	log.Printf("Load pattern: %s", loadPatternStr)

	// Simulate traffic
	startTime := time.Now()
	var completedRequests uint64
	for n := 0; n < b.N; n++ {
		switch loadPattern {
		case 0: // Round-robin
			worker := workers[n%numWorkers]
			worker.Requests <- struct{}{}
			wg.Add(1)
		case 1: // Random
			worker := workers[rand.Intn(numWorkers)]
			worker.Requests <- struct{}{}
			wg.Add(1)
		}

		atomic.AddUint64(&completedRequests, 1)
	}

	// Wait for all workers to finish
	wg.Wait()

	// Calculate duration
	duration := time.Since(startTime)

	// Log metrics
	log.Printf("Benchmark completed: %s", b.Name())
	log.Printf("Benchmark duration: %v", duration)
	log.Printf("Total requests: %d", completedRequests)

	if duration > 0 {
		opsPerSecond := float64(completedRequests) / duration.Seconds()
		log.Printf("Operations per second: %.2f", opsPerSecond)
	} else {