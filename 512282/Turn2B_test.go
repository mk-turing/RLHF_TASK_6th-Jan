package _12282

import (
	"math/rand"
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

// Sliding window analysis with inlined functions
func slidingWindowAnalysisInlined(data []int, windowSize int, funcF func(int) int) []int {
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

// Non-inlined function used in the sliding window analysis
func nonInlinedFunction(x int) int {
	return x * x
}

// Inlined function 1: Simple addition
func inlinedFunction1(x int) int {
	return x + x
}

// Inlined function 2: Simple multiplication
func inlinedFunction2(x int) int {
	return x * x
}

// Inlined function 3: Complex calculation (sin(x) + cos(x))
func inlinedFunction3(x int) int {
	// Using math package is not possible in benchmark, using approximation instead
	const pi = 3.14
	return int((0.5 * (1 + x%2)) * pi * float64(x))
}

// Inlined function 4: Using local variable
func inlinedFunction4(x int) int {
	temp := x * x
	return temp
}

// Sliding window analysis with non-inlined functions
func slidingWindowAnalysisNonInlined(data []int, windowSize int, funcF func(int) int) []int {
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

func BenchmarkSlidingWindowAnalysisInlined(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := nonInlinedFunction

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined1(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlinedFunction1
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined2(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlinedFunction2
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined3(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlinedFunction3
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}
func BenchmarkSlidingWindowAnalysisInlined4(b *testing.B) {