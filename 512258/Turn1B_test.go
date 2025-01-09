package _12258

import (
	"runtime"
	"testing"
)

func MyFunctionToBenchmark() []int {
	// Your function implementation here
	// For this example, let's just allocate a large slice of memory
	var data []int
	data = make([]int, 1000000) // Allocate 1 MB of memory
	return data
	// Do some work with the data (if needed)
}

func TestMyFunctionToBenchmark(t *testing.T) {
	for i := 0; i < 10; i++ {
		MyFunctionToBenchmark()
		// Force garbage collection to get accurate memory usage after each iteration
		runtime.GC()
	}
}

func BenchmarkMyFunctionToBenchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MyFunctionToBenchmark()
	}
}
