package _12258

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// Example functions to benchmark with varying input sizes
func exampleFunction1(size int) {
	// Simulate memory usage and CPU computation based on input size
	data := make([]int, size)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Int()
	}
	for i := 0; i < 100000; i++ {
		// Simulate CPU computation
		for j := 0; j < len(data); j++ {
			data[j] = data[j] * 2
		}
	}
}

func exampleFunction2(size int) {
	// Simulate memory usage and CPU computation differently
	data := make([]int, size)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Int()
	}
	for i := 0; i < 100000; i++ {
		// Simulate CPU computation with a different operation
		for j := 0; j < len(data); j++ {
			data[j] = data[j] ^ 1
		}
	}
}

func BenchmarkExampleFunctions(b *testing.B) {
	inputSizes := []int{1000, 10000, 100000, 1000000}
	for _, size := range inputSizes {
		b.Run(fmt.Sprintf("ExampleFunction1-%d", size), func(b *testing.B) {
			benchmarkFunction(b, exampleFunction1, size)
		})
		b.Run(fmt.Sprintf("ExampleFunction2-%d", size), func(b *testing.B) {
			benchmarkFunction(b, exampleFunction2, size)
		})
	}
}

func benchmarkFunction(b *testing.B, f func(int), size int) {
	// Get the initial memory stats
	var initialStats runtime.MemStats
	runtime.ReadMemStats(&initialStats)

	// Run the benchmark
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		f(size)
	}
	b.StopTimer()

	// Get the final memory stats
	var finalStats runtime.MemStats
	runtime.ReadMemStats(&finalStats)

	// Calculate memory usage
	memoryUsage := finalStats.Alloc - initialStats.Alloc
	b.SetBytes(int64(memoryUsage))

	// Calculate the number of memory allocations
	allocs := finalStats.Mallocs - initialStats.Mallocs

	// Output results
	fmt.Printf("%s      %d ns/op      %.2f MB/op      %.2f allocs/op\n",
		b.Name(), b.N, float64(memoryUsage)/1024/1024, float64(allocs))
}

func getCPUUsage() (float64, error) {
	// Command to get CPU usage (system-dependent)
	cmd := exec.Command("top", "-b", "-n", "1", "-p", fmt.Sprint(os.Getpid()))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	// ... (Same as before)
}

func parsePercentage(str string) float64 {
	// ... (Same as before)
}
