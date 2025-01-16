package _12259

import (
	"fmt"
	"math/rand"
	"runtime"
	"runtime/debug"
	"testing"
)

// BubbleSort function to be benchmarked
func BubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}

// generateRandomArray generates a random array of given size
func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(size)
	}
	return arr
}

func BenchmarkBubbleSortWithMemory(b *testing.B) {
	// Different array sizes
	sizes := []int{1000, 10000, 100000}
	// Different CPU configurations
	cpus := []int{1, 2, runtime.NumCPU()}
	// Different memory configurations (in MB)
	memConfs := []int{16, 32, 64, 128, 256}

	for _, cpusToUse := range cpus {
		runtime.GOMAXPROCS(cpusToUse)
		for _, memConf := range memConfs {
			// Set the memory limit for the current benchmark run
			debug.SetMemoryLimit(int64(memConf) * 1024 * 1024)

			for _, size := range sizes {
				b.Run(fmt.Sprintf("cpus-%d/mem-%dMB/size-%d", cpusToUse, memConf, size), func(b *testing.B) {
					// Reset the allocated memory for each run
					debug.FreeOSMemory()

					// Generate a random array for each benchmark run
					arr := generateRandomArray(size)

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						// Copy the array to avoid modifying the original
						copyArr := make([]int, len(arr))
						copy(copyArr, arr)

						BubbleSort(copyArr)
					}

					// Measure memory allocation after the benchmark
					var m runtime.MemStats
					runtime.ReadMemStats(&m)
					b.ReportMetric(float64(m.Alloc)/float64(1024), "alloc-KB")
				})
			}
		}
	}
}
