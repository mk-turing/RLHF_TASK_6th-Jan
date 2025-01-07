package _93902

import "testing"

func factorial(n int) int {
	if n < 0 {
		return 0
	}
	if n == 0 || n == 1 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		factorial(12)
	}
}
