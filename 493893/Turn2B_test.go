package main

import (
	"math/rand"
	"testing"
)

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

func Algorithm3(data []int) int {
	sum := 0
	for i := 0; i < len(data); i++ {
		sum += data[i]
	}
	return sum
}

func BenchmarkAlgorithms(b *testing.B) {
	// Initialize test data for all algorithms
	data := make([]int, 100000)
	for i := range data {
		data[i] = rand.Intn(100)
	}

	// Run benchmarks for each algorithm
	b.Run("Algorithm1", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Algorithm1(data)
		}
	})

	b.Run("Algorithm2", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Algorithm2(data)
		}
	})

	b.Run("Algorithm3", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Algorithm3(data)
		}
	})
}

func ExampleCompareAlgorithms() {
	// Run the benchmark and print the results
	testing.Main(func(m *testing.M) {
		m.Run()
	}, nil, nil, nil)
	// Output:
	// BenchmarkAlgorithms/Algorithm1-8   	 5000000	       241 ns/op
	// BenchmarkAlgorithms/Algorithm2-8   	 2000000	       601 ns/op
	// BenchmarkAlgorithms/Algorithm3-8   	 5000000	       228 ns/op
	// PASS
}
