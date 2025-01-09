package main

import "testing"

func sumAndDiffInts(a, b int) (int, int) {
	return a + b, a - b
}
func sumAndDiffFloats(a, b float64) (float64, float64) {
	return a + b, a - b
}

// Benchmark functions
func BenchmarkSumAndDiffInts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum, diff := sumAndDiffInts(i, b.N-i)
		_ = sum
		_ = diff
	}
}
func BenchmarkSumAndDiffFloats(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sum, diff := sumAndDiffFloats(float64(i), float64(b.N-i))
		_ = sum
		_ = diff
	}
}

type Result struct {
	Sum  int
	Diff int
}

func sumAndDiffStruct(a, b int) Result {
	return Result{Sum: a + b, Diff: a - b}
}

// Benchmark function
func BenchmarkSumAndDiffStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result := sumAndDiffStruct(i, b.N-i)
		_ = result.Sum
		_ = result.Diff
	}
}
