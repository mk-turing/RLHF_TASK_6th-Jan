package _12260

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

// BenchmarkResult holds the benchmark result data
type BenchmarkResult struct {
	Time       time.Time     `json:"time"`
	Operations int           `json:"operations"`
	Duration   time.Duration `json:"duration"`
	Rate       float64       `json:"rate"`
}

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

	b.StopTimer()
	b.Logf("Total sum: %d", totalSum)
}

func recordBenchmarkResult(result BenchmarkResult) {
	// Open the results file in append mode
	file, err := os.OpenFile("benchmark_results.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening results file: %v", err)
	}
	defer file.Close()

	// Encode the benchmark result to JSON and write it to the file
	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding benchmark result: %v", err)
	}
	fmt.Fprintln(file, string(resultJSON))
}

func main() {
	// Run the benchmark immediately
	b := testing.Benchmark(BenchmarkAddConcurrently)

	// Calculate benchmark result
	result := BenchmarkResult{
		Time:       time.Now(),
		Operations: b.N,
		Duration:   b.Duration,
		Rate:       float64(b.N) / b.Duration.Seconds(),
	}

	// Record the benchmark result
	recordBenchmarkResult(result)

	// Optionally, you can schedule the benchmark to run at regular intervals using a cron job
	c := time.NewTicker(1 * time.Hour) // Run every hour
	defer c.Stop()

	for range c.C {
		// Run the benchmark again and record the result
		b := testing.Benchmark(BenchmarkAddConcurrently)
		result := BenchmarkResult{
			Time:       time.Now(),
			Operations: b.N,
			Duration:   b.Duration,
			Rate:       float64(b.N) / b.Duration.Seconds(),
		}
		recordBenchmarkResult(result)
	}
}
