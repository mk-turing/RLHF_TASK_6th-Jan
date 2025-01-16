package _12259

import (
	"math/rand"
	"runtime"
	"testing"
)

func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	pivot := arr[len(arr)/2]
	left := []int{}
	right := []int{}
	for _, num := range arr {
		if num < pivot {
			left = append(left, num)
		} else if num > pivot {
			right = append(right, num)
		}
	}
	return append(QuickSort(left), pivot, QuickSort(right)...)
}

// BenchmarkQuickSort measures the execution time and memory usage of QuickSort
func BenchmarkQuickSort(b *testing.B) {
	// Define the size of the array to sort
	const arraySize = 100000
	// Create a random array
	arr := make([]int, arraySize)
	for i := range arr {
		arr[i] = int(rand.Int63())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Start memory profiling
		runtime.GC()
		memBefore := runtime.ReadMemStats()

		// Sort the array
		QuickSort(arr)

		// End memory profiling
		memAfter := runtime.ReadMemStats()

		// Calculate memory usage
		memoryUsage := memAfter.Alloc - memBefore.Alloc

		// Record the memory usage
		b.SetBytes(int64(memoryUsage))
	}
}
