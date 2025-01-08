package main

import (
	"testing"
)

func sumRangeManual(start, end int) int {
	total := 0
	for i := start; i <= end; i++ {
		total += i
	}
	return total
}

func BenchmarkSumRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sumRangeManual(1, 1000000)
	}
}

func TestMain(m *testing.M) {
	testing.Init()
	// Use testing.M.Run to run the benchmarks and tests together
	m.Run()
}
