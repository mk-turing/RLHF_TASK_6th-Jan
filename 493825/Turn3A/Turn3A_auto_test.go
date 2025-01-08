package main

import (
	"runtime"
	"sync"
	"testing"
)

func sumRangeManual(start, end int) int {
	total := 0
	for i := start; i <= end; i++ {
		total += i
	}
	return total
}

// Automated concurrent benchmark for large range
func BenchmarkSumRangeLargeConcurrent(b *testing.B) {
	rangeSize := 1000000 / b.N

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		var results []int

		// Spawn goroutines for each CPU core
		for j := 0; j < runtime.NumCPU(); j++ {
			wg.Add(1)
			start := i*rangeSize + j*rangeSize/b.N
			end := i*rangeSize + (j+1)*rangeSize/b.N
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

		if total != 1000000*b.N*(b.N+1)/2 {
			b.Fatalf("computation error: expected %d, got %d", 1000000*b.N*(b.N+1)/2, total)
		}
	}
}

func TestMain(m *testing.M) {
	testing.Init()
	// Use testing.M.Run to run the benchmarks and tests together
	m.Run()
}
