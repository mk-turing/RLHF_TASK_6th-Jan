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

func benchmarkManual() {
	runtime.GOMAXPROCS(1) // Ensure a single CPU core is used for consistency

	start := time.Now()
	sumRangeManual(1, 1000000)
	end := time.Now()

	fmt.Printf("Manual Benchmark: %.6f seconds\n", end.Sub(start).Seconds())
}

func main() {
	benchmarkManual()
}
