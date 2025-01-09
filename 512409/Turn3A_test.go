package main

import (
	"sync"
	"testing"
)

// Function returning multiple values
func sumAndDiff(a, b int) (sum int, diff int) {
	sum = a + b
	diff = a - b
	return
}

// Function returning multiple values in a struct
type Result struct {
	Sum  int
	Diff int
}

func sumAndDiffStruct(a, b int) Result {
	return Result{Sum: a + b, Diff: a - b}
}

// Benchmark function for sequential execution
func BenchmarkSumAndDiffSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum, diff := sumAndDiff(i, b.N-i)
		_ = sum
		_ = diff
	}
}

// Benchmark function for concurrent execution using goroutines
func BenchmarkSumAndDiffConcurrent(b *testing.B) {
	const numGoroutines = 10
	var wg sync.WaitGroup

	b.SetParallelism(numGoroutines)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sum, diff := sumAndDiff(i, b.N-i)
			_ = sum
			_ = diff
		}(i)
	}

	wg.Wait()
}

// Benchmark function for sequential execution using struct
func BenchmarkSumAndDiffStructSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := sumAndDiffStruct(i, b.N-i)
		_ = result.Sum
		_ = result.Diff
	}
}

// Benchmark function for concurrent execution using struct
func BenchmarkSumAndDiffStructConcurrent(b *testing.B) {
	const numGoroutines = 10
	var wg sync.WaitGroup

	b.SetParallelism(numGoroutines)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			result := sumAndDiffStruct(i, b.N-i)
			_ = result.Sum
			_ = result.Diff
		}(i)
	}

	wg.Wait()
}
