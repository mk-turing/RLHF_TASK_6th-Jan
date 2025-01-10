package main

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

// TestMemoryAllocation tests memory allocation and returns the result.
func TestMemoryAllocation(t *testing.T) {
	var result int
	for i := 0; i < 1000000; i++ {
		slice := make([]int, 0, 1000) // Preallocate capacity
		for j := 0; j < 1000; j++ {
			slice = append(slice, j)
		}
		result = sumSlice(slice) // Calculate sum (or just the first element)
	}
	t.Logf("Sum of first element in slice: %d", result)
}

// sumSlice returns the first element of the slice or 0 if the slice is empty.
func sumSlice(slice []int) int {
	if len(slice) == 0 { // Ensure slice is not empty
		return 0
	}
	return slice[0]
}

// BenchmarkMemoryAllocation benchmarks memory allocation performance.
func BenchmarkMemoryAllocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var result int
		for i := 0; i < 1000000; i++ {
			slice := make([]int, 0, 1000) // Preallocate capacity
			for j := 0; j < 1000; j++ {
				slice = append(slice, j)
			}
			result = sumSlice(slice) // Calculate sum (or just the first element)
		}
		b.Logf("Sum of first element in slice: %d", result)
	}
}

// main function to demonstrate memory statistics and trigger benchmark.
func main() {
	fmt.Println("Memory information before tests:")
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	fmt.Printf("Total Alloc: %d bytes\n", memStats.TotalAlloc)
	fmt.Printf("Sys: %d bytes\n", memStats.Sys)
	fmt.Printf("Num GC: %d\n", memStats.NumGC)
	fmt.Printf("Pause Total Nanoseconds: %d ns\n", memStats.PauseTotalNs)

	// Enable CPU profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create CPU profile: %v\n", err)
		os.Exit(2)
	}
	runtime.GC()
	runtime.PPROF_CPU = f
	defer runtime.PPROF_CPU.Close()

	// Enable memory profiling
	f, err = os.Create("mem.prof")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create memory profile: %v\n", err)
		os.Exit(2)
	}
	defer f.Close()
	runtime.ReadMemStats(memStats)
	runtime.PPROF_MEMORY = f
	runtime.GC()
	defer runtime.PPROF_MEMORY.Close()

	// Run tests and benchmarks
	testing.Run(os.Args[0], testing.MainStart)

	// Print memory stats after tests.
	runtime.GC()
	runtime.ReadMemStats(memStats)
	fmt.Println("\nMemory information after GC:")
	fmt.Printf("Total Alloc: %d bytes\n", memStats.TotalAlloc)
	fmt.Printf("Sys: %d bytes\n", memStats.Sys)
	fmt.Printf("Num GC: %d\n", memStats.NumGC)
	fmt.Printf("Pause Total Nanoseconds: %d ns\n", memStats.PauseTotalNs)
}
