package _12258

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

var (
	inputSizes = []int{10000, 100000, 1000000} // Different input sizes to test
)

// Function 1: Allocates slices in a loop
func func1(size int) {
	for i := 0; i < size; i++ {
		_ = make([]int, 100)
	}
}

// Function 2: Uses make with a known capacity in each iteration
func func2(size int) {
	var s []int
	for i := 0; i < size; i++ {
		s = make([]int, 0, 100)
		for j := 0; j < 100; j++ {
			s = append(s, j)
		}
	}
}

// Function 3: Uses append to grow a slice in each iteration
func func3(size int) {
	var s []int
	for i := 0; i < size; i++ {
		for j := 0; j < 100; j++ {
			s = append(s, j)
		}
		s = s[:0] // Clear the slice to avoid reallocations
	}
}

func benchmarkFunction(b *testing.B, f func(int), size int) {
	// Get the initial memory stats
	var initialStats runtime.MemStats
	runtime.ReadMemStats(&initialStats)

	// Run the benchmark
	b.StartTimer() // Start timing only the benchmarked code
	for n := 0; n < b.N; n++ {
		f(size)
	}
	b.StopTimer() // Stop timing after the benchmark

	// Get the final memory stats
	var finalStats runtime.MemStats
	runtime.ReadMemStats(&finalStats)

	// Calculate memory usage
	memoryUsage := finalStats.Alloc - initialStats.Alloc
	b.SetBytes(int64(memoryUsage))

	// Calculate the number of memory allocations
	allocs := finalStats.Mallocs - initialStats.Mallocs

	// Output results
	fmt.Printf("Benchmark%s-%d      %d ns/op      %.2f MB/op      %.2f allocs/op\n",
		strings.ReplaceAll(fmt.Sprintf("%T", f), "func ", ""), size, b.N,
		float64(memoryUsage)/1024/1024, float64(allocs))
}

// A more reliable method for getting CPU usage in Go
func getCPUUsage() (float64, error) {
	// Command to get CPU usage (system-dependent)
	cmd := exec.Command("top", "-b", "-n", "1", "-p", fmt.Sprint(os.Getpid()))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// Parse CPU usage from the output (system-dependent)
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "Cpu(s):") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				user := fields[1]
				return parsePercentage(user), nil
			}
		}
	}

	return 0, fmt.Errorf("cannot find CPU usage in top output")
}

func parsePercentage(str string) float64 {
	// Convert the string percentage to float before multiplying by 100
	percentage, err := strconv.ParseFloat(strings.TrimRight(str, "%"), 64)
	if err != nil {
		return 0
	}
	return percentage
}

func BenchmarkTest(t *testing.B) {
	for _, size := range inputSizes {
		fmt.Printf("Testing size: %d\n", size)
		benchmarkFunction(t, func1, size)
		benchmarkFunction(t, func2, size)
		benchmarkFunction(t, func3, size)
	}
}
