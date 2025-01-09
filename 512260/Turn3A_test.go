package _12260

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

// BenchmarkResult represents a single benchmark result
type BenchmarkResult struct {
	Timestamp     time.Time
	RunTimeNs     int64
	TotalSum      int
	NumGoroutines int
}

func add(a, b int) int {
	return a + b
}

// BenchmarkAddConcurrently performs a concurrent benchmark and records the result
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

	// Start the timer
	b.StartTimer()

	// Wait for all goroutines to finish
	wg.Wait()

	// Stop the timer
	b.StopTimer()

	// Record the result
	result := BenchmarkResult{
		Timestamp:     time.Now(),
		TotalSum:      totalSum,
		RunTimeNs:     b.Elapsed().Nanoseconds(),
		NumGoroutines: numGoroutines,
	}

	// Save the result to a file
	saveResult(result)
}

// saveResult writes the benchmark result to a file
func saveResult(result BenchmarkResult) {
	filename := "benchmark_results.json"
	resultsFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening results file:", err)
		return
	}
	defer resultsFile.Close()

	resultsFile.WriteString(result.Timestamp.Format(time.RFC3339) + ",\n")
	resultsFile.WriteString(fmt.Sprintf("{%v},\n", result))
}

func main() {
	// Run the benchmark
	testing.Benchmark(BenchmarkAddConcurrently)

	// Print help for review
	fmt.Println("Benchmark results have been saved to benchmark_results.json")
}
