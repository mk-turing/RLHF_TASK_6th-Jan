package _12258

import (
	"runtime"
	"testing"
	"time"
)

// Example function to benchmark
func exampleFunction() {
	// Simulate memory usage and CPU work
	for i := 0; i < 100000; i++ {
		// Allocate a slice
		_ = make([]int, 1000)
		// Simulate some CPU work
		for j := 0; j < 10000; j++ {
			_ = j * j
		}
	}
}

func BenchmarkExampleFunction(b *testing.B) {
	// Reset CPU usage counter
	runtime.GOMAXPROCS(runtime.NumCPU())
	var cpuTime time.Duration

	// Get the initial memory stats
	var initialStats runtime.MemStats
	runtime.ReadMemStats(&initialStats)

	b.ResetTimer()
	// Run the benchmark
	for n := 0; n < b.N; n++ {
		// Measure CPU usage
		start := time.Now()
		exampleFunction()
		cpuTime += time.Since(start)
	}

	// Get the final memory stats
	var finalStats runtime.MemStats
	runtime.ReadMemStats(&finalStats)

	// Calculate memory usage
	memoryUsage := finalStats.Alloc - initialStats.Alloc
	b.SetBytes(int64(memoryUsage))

	// Report CPU usage
	b.SetReport(&testing.BenchmarkReport{
		MemoryBytes: memoryUsage,
		AllocsPerOp: uint64(finalStats.Allocs - initialStats.Allocs),
		CPU:         cpuTime,
	})
}

func main() {
	testing.MainStart(testing.TestMainStart{}, nil, nil, nil, false, "")
}
