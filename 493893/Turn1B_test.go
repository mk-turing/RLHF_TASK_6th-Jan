package main

import (
	"testing"
)

func Algorithm1(data []int) int {
	sum := 0
	for _, num := range data {
		sum += num
	}
	return sum
}

// Algorithm2 multiplies all the numbers in the input slice
func Algorithm2(data []int) int {
	product := 1
	for _, num := range data {
		product *= num
	}
	return product
}

func BenchmarkAlgorithm1(b *testing.B) {
	// Initialize test data for the algorithm
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Your algorithm to be benchmarked here
		Algorithm1(data)
	}
}

func BenchmarkAlgorithm2(b *testing.B) {
	// Initialize test data for the algorithm
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Your other algorithm to be benchmarked here
		Algorithm2(data)
	}
}
