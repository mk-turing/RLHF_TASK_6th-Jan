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

// Example function to benchmark
func exampleFunction() {
	// Simulate memory usage and CPU computation
	for i := 0; i < 100000; i++ {
		// Allocate a slice (this will be garbage collected eventually)
		_ = make([]int, 1000)
		// Simulate CPU computation
		for j := 0; j < 10000; j++ {
			_ = i + j
		}
	}
}

func BenchmarkExampleFunction(b *testing.B) {
	// Get the initial memory stats
	var initialStats runtime.MemStats
	runtime.ReadMemStats(&initialStats)

	// Run the benchmark
	b.StartTimer() // Start timing only the benchmarked code
	for n := 0; n < b.N; n++ {
		exampleFunction()
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
	fmt.Printf("BenchmarkExampleFunction-4      %d ns/op      %.2f MB/op      %.2f allocs/op\n",
		b.N, float64(memoryUsage)/1024/1024, float64(allocs))
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
