package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func sumRangeManual(start, end int) int {
	total := 0
	for i := start; i <= end; i++ {
		total += i
	}
	return total
}

// Concurrent benchmark for large range using goroutines
func benchmarkLargeRangeConcurrentManual() {
	runtime.GOMAXPROCS(1) // Ensure a single CPU core is used for consistency

	var wg sync.WaitGroup
	var results []int
	rangeSize := 1000000 / runtime.NumCPU()

	// Spawn goroutines for each CPU core
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		start := i*rangeSize + 1
		end := (i + 1) * rangeSize
		go func() {
			defer wg.Done()
			results = append(results, sumRangeManual(start, end))
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Sum up all partial results
	total := 0
	for _, result := range results {
		total += result
	}

	// Measure total execution time
	start := time.Now()
	sumRangeManual(1, 1000000)
	end := time.Now()

	fmt.Printf("Large Range (1-1000000) Concurrent Manual Benchmark: %.6f seconds\n", end.Sub(start).Seconds())
}

func main() {
	benchmarkLargeRangeConcurrentManual()
}
