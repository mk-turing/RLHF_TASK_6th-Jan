package _12282

import (
	"testing"
)

// Non-inlined sliding window function
func sumSlidingWindowNonInlined(nums []int, k int) []int {
	if len(nums) == 0 || k <= 0 {
		return nil
	}

	var result []int
	windowSum := sumHelper(nums[:k])
	result = append(result, windowSum)

	for i := k; i < len(nums); i++ {
		windowSum += nums[i] - nums[i-k]
		result = append(result, windowSum)
	}

	return result
}

// Helper function for non-inlined version to compute the sum
func sumHelper(arr []int) int {
	s := 0
	for _, num := range arr {
		s += num
	}
	return s
}

// Inlined sliding window function
func sumSlidingWindowInlined(nums []int, k int) []int {
	if len(nums) == 0 || k <= 0 {
		return nil
	}

	var result []int
	s := 0
	for i := 0; i < k; i++ {
		s += nums[i]
	}
	result = append(result, s)

	for i := k; i < len(nums); i++ {
		s += nums[i] - nums[i-k]
		result = append(result, s)
	}

	return result
}

func BenchmarkSumSlidingWindowNonInlined(b *testing.B) {
	nums := make([]int, 100000)
	for i := range nums {
		nums[i] = i
	}
	k := 1000
	for i := 0; i < b.N; i++ {
		sumSlidingWindowNonInlined(nums, k)
	}
}

func BenchmarkSumSlidingWindowInlined(b *testing.B) {
	nums := make([]int, 100000)
	for i := range nums {
		nums[i] = i
	}
	k := 1000
	for i := 0; i < b.N; i++ {
		sumSlidingWindowInlined(nums, k)
	}
}
