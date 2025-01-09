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

// Inlined function variation 1: Squaring
func inlineSquaring(x int) int {
	return x * x
}

// Inlined function variation 2: Cubing
func inlineCubing(x int) int {
	return x * x * x
}

// Inlined function variation 3: Adding 1
func inlineAddingOne(x int) int {
	return x + 1
}

// Inlined function variation 4: Multiplying by 2
func inlineMultiplyingByTwo(x int) int {
	return x * 2
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

func BenchmarkSlidingWindowAnalysisInlined_Squaring(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlineSquaring

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined_Cubing(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlineCubing

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined_AddingOne(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlineAddingOne

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}

func BenchmarkSlidingWindowAnalysisInlined_MultiplyingByTwo(b *testing.B) {
	data := generateRandomData(100000, 100)
	windowSize := 100
	funcF := inlineMultiplyingByTwo

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slidingWindowAnalysisInlined(data, windowSize, funcF)
	}
}
