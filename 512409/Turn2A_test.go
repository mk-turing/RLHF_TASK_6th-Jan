package main

import (
	"testing"
)

// Function returning a single primitive type
func sumPrimitive(a, b int) int {
	return a + b
}

// Function returning multiple primitive types
func sumMultiplePrimitive(a, b int) (sum int, overflow bool) {
	sum, overflow = a+b, false
	return
}

// Function returning a single struct
type Result struct {
	Sum      int
	Overflow bool
}

func sumStruct(a, b int) Result {
	return Result{Sum: a + b, Overflow: false}
}

// Function returning a slice (heap-allocated)
func sliceBenchmark() []int {
	return make([]int, 10000)
}

func TestBenchmark(t *testing.T) {
	testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sumPrimitive(i, b.N-i)
		}
	})

	testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sumMultiplePrimitive(i, b.N-i)
		}
	})

	testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sumStruct(i, b.N-i)
		}
	})

	testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sliceBenchmark()
		}
	})
}
