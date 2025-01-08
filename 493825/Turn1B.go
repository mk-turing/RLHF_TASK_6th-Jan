package main

import (
	"fmt"
	"github.com/pkg/profile"
	"testing"
	"time"
)

// Function to be benchmarked
func sum(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

// Manual benchmark function
func manualBenchmark(n int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum(n)
	}
}

// Automated benchmark function
func automatedBenchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, n := range []int{100, 1000, 10000, 100000} {
			sum(n)
		}
	}
}
func main() {
	// Manual benchmark
	n := 100000
	b := &testing.B{}
	b.ResetTimer()
	manualBenchmark(n, b)
	b.StopTimer()
	fmt.Printf("Manual benchmark for n=%d: %v\n", n, b.Elapsed())
	// Automated benchmark
	testing.Benchmark(automatedBenchmark)
	// Profile the program to analyze CPU usage
	f := profile.Start(profile.CPUProfile)
	defer f.Stop()
	time.Sleep(5 * time.Second)
}
