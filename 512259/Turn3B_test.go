package _12259

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
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

// BenchmarkBubbleSortWithDiskUsage benchmarks the BubbleSort function with different array sizes, CPU configurations, and disk speeds.
func BenchmarkBubbleSortWithDiskUsage(b *testing.B) {
	// Different array sizes
	sizes := []int{10000, 100000}
	// Different CPU configurations
	cpus := []int{1, 2, runtime.NumCPU()}
	// Different disk speeds (in MB/s)
	diskSpeeds := []int{1, 10, 100, 500}

	// To monitor disk usage
	var wg sync.WaitGroup

	for _, cpusToUse := range cpus {
		runtime.GOMAXPROCS(cpusToUse)

		for _, diskSpeed := range diskSpeeds {
			wg.Add(1)
			go func(diskSpeed int) {
				defer wg.Done()
				// Create a temporary file for monitoring disk usage
				file, err := ioutil.TempFile("", "bubblesort_test")
				if err != nil {
					b.Fatal(err)
				}
				defer os.Remove(file.Name())

				// Write data to the file at the given disk speed
				go func() {
					for {
						data := make([]byte, 1024*1024) // 1MB of data
						_, _ = file.Write(data)
						// Simulate disk I/O delay based on disk speed
						runtime.Gosched()
					}
				}()

				for _, size := range sizes {
					b.Run(fmt.Sprintf("cpus-%d/disk-%dMBps/size-%d", cpusToUse, diskSpeed, size), func(b *testing.B) {
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
					})
				}

				// Stop writing data to the file after the benchmark
				close(done)
			}(diskSpeed)
		}
	}
	wg.Wait()
}
