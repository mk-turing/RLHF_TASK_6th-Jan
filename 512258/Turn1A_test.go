package _12258

import (
	"runtime"
	"testing"
)

// Example function to benchmark
func exampleFunction() {
	// Simulate memory usage
	for i := 0; i < 100000; i++ {
		// Allocate a slice (this will be garbage collected eventually)
		_ = make([]int, 1000)
	}
}

func BenchmarkExampleFunction(b *testing.B) {
	// Get the initial memory stats
	var initialStats runtime.MemStats
	runtime.ReadMemStats(&initialStats)

	// Run the benchmark
	for n := 0; n < b.N; n++ {
		exampleFunction()
	}

	// Get the final memory stats
	var finalStats runtime.MemStats
	runtime.ReadMemStats(&finalStats)

	// Calculate memory usage
	memoryUsage := finalStats.Alloc - initialStats.Alloc
	b.SetBytes(int64(memoryUsage))
}
