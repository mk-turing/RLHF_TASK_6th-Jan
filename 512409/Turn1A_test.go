package main

import (
	"testing"
)

// Single return value
func sumSingle(a, b int) int {
	return a + b
}

// Multiple return values
func sumMultiple(a, b int) (sum int, overflow bool) {
	sum, overflow = a+b, false
	return
}

func TestBenchmark(t *testing.T) {
	testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sumSingle(i, b.N-i)
		}
	})

	testing.Benchmark(func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sumMultiple(i, b.N-i)
		}
	})
}
