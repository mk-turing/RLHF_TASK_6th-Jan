package _12282

import (
	"math/rand"
	"runtime"
	"testing"
)

// Generate a slice of random integers with a given size and maximum value
func generateRandomData(size int, maxValue int) []int {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(maxValue)
	}
	return data
}

// Simple inlined function
func simpleInline(x int) int {
	return x + 1
}

// Complex inlined function (with multiple nested loops)
func complexInline(x int) int {
	sum := 0
	for i := 0; i < x; i++ {
		for j := 0; j < i; j++ {
			sum += i + j
		}
	}
	return sum
}

// Non-inlined function
func nonInlinedFunction(x int) int {
	return x * x
}

// Sliding window analysis with inlined functions
func slidingWindowAnalysis(data []int, windowSize int, funcF func(int) int) []int {
	result := make([]int, len(data)-windowSize+1)
	for i := 0; i < len(result); i++ {
		sum := 0
		for j := i; j < i+windowSize; j++ {
			sum += funcF(data[j])
		}
		result[i] = sum
	}
	return result
}

// Benchmark with memory usage tracking
func BenchmarkSlidingWindowAnalysis_SimpleInline(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100

	var memStats runtime.MemStats
	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&memStats)
		b.StartTimer()
		slidingWindowAnalysis(data, windowSize, simpleInline)
		b.StopTimer()
		b.RecordMemoryAllocs(memStats.Alloc - memStats.Free)
	}
}

func BenchmarkSlidingWindowAnalysis_ComplexInline(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100

	var memStats runtime.MemStats
	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&memStats)
		b.StartTimer()
		slidingWindowAnalysis(data, windowSize, complexInline)
		b.StopTimer()
		b.RecordMemoryAllocs(memStats.Alloc - memStats.Free)
	}
}

func BenchmarkSlidingWindowAnalysis_NonInlined(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100

	var memStats runtime.MemStats
	for i := 0; i < b.N; i++ {
		runtime.ReadMemStats(&memStats)
		b.StartTimer()
		slidingWindowAnalysis(data, windowSize, nonInlinedFunction)
		b.StopTimer()
		b.RecordMemoryAllocs(memStats.Alloc - memStats.Free)
	}
}
