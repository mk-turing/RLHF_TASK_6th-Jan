package main

import (
	"testing"
)

// Define functions to be benchmarked
func Algorithm1(data []int) int {
	sum := 0
	for _, num := range data {
		sum += num
	}
	return sum
}

func Algorithm2(data []int) int {
	product := 1
	for _, num := range data {
		product *= num
	}
	return product
}

// Ensure functions are unique to prevent name collisions
func Algorithm3(data []int) float64 {
	var sum float64 = 0
	for _, num := range data {
		sum += float64(num)
	}
	return sum / float64(len(data))
}

// Benchmark functions
func BenchmarkAlgorithm1(b *testing.B) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Algorithm1(data)
	}
}

func BenchmarkAlgorithm2(b *testing.B) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Algorithm2(data)
	}
}

func BenchmarkAlgorithm3(b *testing.B) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Algorithm3(data)
	}
}
