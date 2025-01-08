package main

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func sumRangeManual(start, end int) int {
	total := 0
	for i := start; i <= end; i++ {
		total += i
	}
	return total
}

// Custom benchmark functions for different edge cases and input sizes
func BenchmarkSumRangeSmall(b *testing.B) {
	benchmarkSumRange(b, 1, 10)
}

func BenchmarkSumRangeMedium(b *testing.B) {
	benchmarkSumRange(b, 1, 1000)
}

func BenchmarkSumRangeLarge(b *testing.B) {
	benchmarkSumRange(b, 1, 1000000)
}

func BenchmarkSumRangeNegative(b *testing.B) {
	benchmarkSumRange(b, -1, 1000)
}

func BenchmarkSumRangeWithLargeStart(b *testing.B) {
	benchmarkSumRange(b, 100000000, 100000010)
}

func benchmarkSumRange(b *testing.B, start, end int) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sumRangeManual(start, end)
	}
	b.StopTimer()
}

func main() {
	// Run the standard automated benchmarking tool (go test -bench .)
	testing.MainStart(testing.TestMainStartOptions{})

	// Manual benchmarking for comparison
	fmt.Println("--- Manual Benchmarks ---")
	benchmarkManual(1, 10)
	benchmarkManual(1, 1000)
	benchmarkManual(1, 1000000)
	benchmarkManual(-1, 1000)
	benchmarkManual(100000000, 100000010)
}

func benchmarkManual(start, end int) {
	runtime.GOMAXPROCS(1)
	startTime := time.Now()
	for i := 0; i < 10000; i++ {
		sumRangeManual(start, end)
	}
	duration := time.Since(startTime)
	fmt.Printf("Manual Benchmark for [%d, %d]: %.6f seconds (10000 iterations)\n", start, end, duration.Seconds())
}
