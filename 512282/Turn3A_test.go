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

// Non-inlined function used in the sliding window analysis
func nonInlinedFunction(x int) int {
	return x * x
}

// Inlined function variation 1: Simple add-one
func inlineAddOne(x int) int {
	return x + 1
}

// Inlined function variation 2: Simple multiplication
func inlineMultiplyByTwo(x int) int {
	return x * 2
}

// Inlined function variation 3: More complex - square root
func inlineSqrt(x int) int {
	// Fails for negative numbers, suit for greater than 0 integers
	sqrt := 0
	for sqrt*sqrt <= x {
		sqrt++
	}
	return sqrt - 1
}

// Inlined function variation 4: Complex - Factorial with a small result
func inlineFactorial(x int) int {
	factorial := 1
	if x > 1 {
		for i := 2; i <= x; i++ {
			factorial *= i
		}
	}
	return factorial
}

func BenchmarkSlidingWindowAnalysisNonInlined(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := nonInlinedFunction

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisNonInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined_AddOne(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlineAddOne

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined_MultiplyByTwo(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlineMultiplyByTwo

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined_Sqrt(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlineSqrt

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined_Factorial(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlineFactorial

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}
