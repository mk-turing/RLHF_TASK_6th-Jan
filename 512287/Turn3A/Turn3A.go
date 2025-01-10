package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"testing"
	"time"
)

// slicePool contains a queue of reusable slice buffers.
var slicePool = []*[]int{}

// Constants for benchmarking
const (
	NumIterations   = 10
	ProfileInterval = 100 * time.Millisecond
)

func getSlice() *[]int {
	if len(slicePool) > 0 {
		// Manually pop from the slice pool
		slice := slicePool[len(slicePool)-1]
		slicePool = slicePool[:len(slicePool)-1] // Remove the last element
		return slice
	}
	return &[]int{}
}

func returnSlice(slice *[]int) {
	slicePool = append(slicePool, slice)
}

func TestMemoryAllocationOptimized(t *testing.T) {
	var result int
	for i := 0; i < 1000000; i++ {
		slice := getSlice() // Fetch from the pool
		for j := 0; j < 1000; j++ {
			*slice = append(*slice, j)
		}
		result = sumSlice(*slice)
		returnSlice(slice) // Return to the pool
	}
	t.Logf("Sum of first element in slice: %d", result)
}

func BenchmarkMemoryAllocationOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var result int
		for i := 0; i < 1000000; i++ {
			slice := getSlice() // Fetch from the pool
			for j := 0; j < 1000; j++ {
				*slice = append(*slice, j)
			}
			result = sumSlice(*slice)
			returnSlice(slice) // Return to the pool
		}
		b.Logf("Sum of first element in slice: %d", result)
	}
}

func sumSlice(slice []int) int {
	if len(slice) == 0 {
		return 0
	}
	return slice[0]
}

// Optimized memory profiling
func runBenchmarkOptimized(testFunc func(t *testing.T)) {
	// Start CPU profiling
	cpuprofile := fmt.Sprintf("cpuprofile_optimized.%d", time.Now().UnixNano())
	f, err := os.Create(cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Start heap profiling
	heapprofile := fmt.Sprintf("heapprofile_optimized.%d", time.Now().UnixNano())
	heapDump := make(chan bool)
	go func() {
		for range time.Tick(ProfileInterval) {
			// Write heap profile to a file
			f, err := os.Create(heapprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.Lookup("heap").WriteTo(f, 0)
			f.Close()
		}
	}()

	// Run the benchmark multiple times to get stable results
	for i := 0; i < NumIterations; i++ {
		testFunc(&testing.T{})
	}

	// Stop heap profiling
	close(heapDump)
}

func generateReportsOptimized() {
	// Generate CPU profile report
	cmd := exec.Command("go", "tool", "pprof", "cpuprofile_optimized.report", "cpuprofile_optimized.*")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Optimized CPU profile report generated")

	// Generate heap profile report
	cmd = exec.Command("go", "tool", "pprof", "heapprofile_optimized.*", "heapprofile_optimized.report")
	var outHeap bytes.Buffer
	cmd.Stdout = &outHeap
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Optimized heap profile report generated")
	fmt.Println(outHeap.String())
}

func main() {
	// Run optimized benchmark
	runBenchmarkOptimized(TestMemoryAllocationOptimized)

	// Generate optimized reports
	generateReportsOptimized()
}
