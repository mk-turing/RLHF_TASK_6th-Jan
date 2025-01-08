package main

import (
	"fmt"
	"runtime"
	"time"
)

func sumRangeManual(start, end int) int {
	total := 0
	for i := start; i <= end; i++ {
		total += i
	}
	return total
}

// Benchmark small range
func benchmarkSmallRangeManual() {
	runtime.GOMAXPROCS(1) // Ensure a single CPU core is used for consistency

	start := time.Now()
	sumRangeManual(1, 1000)
	end := time.Now()

	fmt.Printf("Small Range (1-1000) Manual Benchmark: %.6f seconds\n", end.Sub(start).Seconds())
}

// Benchmark large range
func benchmarkLargeRangeManual() {
	runtime.GOMAXPROCS(1) // Ensure a single CPU core is used for consistency

	start := time.Now()
	sumRangeManual(1, 1000000)
	end := time.Now()

	fmt.Printf("Large Range (1-1000000) Manual Benchmark: %.6f seconds\n", end.Sub(start).Seconds())
}

// Benchmark edge case: zero range
func benchmarkZeroRangeManual() {
	runtime.GOMAXPROCS(1) // Ensure a single CPU core is used for consistency

	start := time.Now()
	sumRangeManual(1, 1)
	end := time.Now()

	fmt.Printf("Zero Range (1-1) Manual Benchmark: %.6f seconds\n", end.Sub(start).Seconds())
}

func main() {
	benchmarkSmallRangeManual()
	benchmarkLargeRangeManual()
	benchmarkZeroRangeManual()
}
