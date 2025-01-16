package _12259

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"testing"
	"time"
)

const tempDir = "./testdata/bubblesort"

func init() {
	os.MkdirAll(tempDir, 0755)
	rand.Seed(time.Now().UnixNano())
}

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

// createTempFile creates a temporary file with a random array
func createTempFile(size int) string {
	fileName := path.Join(tempDir, fmt.Sprintf("temp-%d.txt", time.Now().UnixNano()))
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	arr := generateRandomArray(size)
	_, err = file.Write([]byte(fmt.Sprint(arr)))
	if err != nil {
		panic(err)
	}

	return fileName
}

// readAndSortFile reads a file, sorts its content, and returns the sorted array
func readAndSortFile(fileName string) []int {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	arr := make([]int, 0, len(bytes))
	for _, b := range bytes {
		arr = append(arr, int(b))
	}

	BubbleSort(arr)

	return arr
}

// BenchmarkBubbleSortWithDiskIO benchmarks BubbleSort with disk I/O operations.
func BenchmarkBubbleSortWithDiskIO(b *testing.B) {
	// Different array sizes
	sizes := []int{1000, 5000, 10000}
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
				fileName := createTempFile(size)
				defer os.Remove(fileName)

				b.Run(fmt.Sprintf("cpus-%d/mem-%dMB/size-%d", cpusToUse, memConf, size), func(b *testing.B) {
					// Reset the allocated memory for each run
					debug.FreeOSMemory()

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						readAndSortFile(fileName)
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
