package _12259

import (
	"fmt"
	"math/rand"
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

// BenchmarkBubbleSort benchmarks the BubbleSort function with different array sizes.
func BenchmarkBubbleSort(b *testing.B) {
	// Define different array sizes to benchmark
	sizes := []int{1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			// Generate a random array for each benchmark run
			arr := generateRandomArray(size)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Copy the array to avoid modifying the original
				copyArr := make([]int, len(arr))
				copy(copyArr, arr)

				BubbleSort(copyArr)
			}
		})
	}
}

// ExampleBubbleSort demonstrates how to use the BubbleSort function.
func ExampleBubbleSort() {
	unsortedArray := []int{64, 34, 25, 12, 22, 11, 90}
	BubbleSort(unsortedArray)
	fmt.Println(unsortedArray) // Output: [11 12 22 25 34 64 90]
}
