package _12258

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func captureMemorySnapshot() *runtime.MemStats {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	return &stats
}

func printMemoryStats(stats *runtime.MemStats) {
	fmt.Printf("Total Alloc: %d\n", stats.TotalAlloc)
	fmt.Printf("Sys: %d\n", stats.Sys)
	fmt.Printf("Heap Alloc: %d\n", stats.HeapAlloc)
	fmt.Printf("Heap Sys: %d\n", stats.HeapSys)
	fmt.Printf("Heap Idle: %d\n", stats.HeapIdle)
	fmt.Printf("Heap Inuse: %d\n", stats.HeapInuse)
	fmt.Printf("Heap Released: %d\n", stats.HeapReleased)
	fmt.Printf("Heap Objects: %d\n", stats.HeapObjects)
	fmt.Printf("Mallocs: %d\n", stats.Mallocs)
	fmt.Printf("Frees: %d\n", stats.Frees)
}

func TestMemoryUsage(t *testing.T) {
	numRuns := 10
	snapshots := make([]*runtime.MemStats, numRuns)

	for i := 0; i < numRuns; i++ {
		// Capture a snapshot before running the function
		before := captureMemorySnapshot()

		// Run the function to test
		testFunction()

		// Capture a snapshot after running the function
		after := captureMemorySnapshot()

		// Calculate the memory usage change
		usageChange := after.HeapAlloc - before.HeapAlloc

		// Print memory stats
		fmt.Printf("Run %d: Memory usage change: %d bytes\n", i+1, usageChange)
		printMemoryStats(after)

		// Store the snapshot for comparison
		snapshots[i] = after

		// Force garbage collection to ensure accurate memory usage
		runtime.GC()

		// Wait for a short period to observe memory trends
		time.Sleep(500 * time.Millisecond)
	}

	// Analyze memory trends
	for i := 1; i < numRuns; i++ {
		if snapshots[i].HeapAlloc > snapshots[i-1].HeapAlloc {
			fmt.Printf("Memory leak detected between runs %d and %d: %d bytes\n", i, i-1, snapshots[i].HeapAlloc-snapshots[i-1].HeapAlloc)
		}
	}
}

func testFunction() {
	// Simulate a memory leak by allocating slices and not freeing them
	for i := 0; i < 100000; i++ {
		_ = make([]int, 1000)
	}
}
