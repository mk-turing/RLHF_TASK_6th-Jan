package _93902

import "testing"

func TestFactorial(t *testing.T) {
	testCases := []struct {
		n   int
		exp int
	}{
		{0, 1},
		{1, 1},
		{3, 6},
		{5, 120},
		{-1, 0},
	}

	for _, tc := range testCases {
		result := factorial(tc.n)
		if result != tc.exp {
			t.Errorf("factorial(%d) = %d, expected %d", tc.n, result, tc.exp)
		}
	}
}

func factorial(n int) int {
	if n < 0 {
		return 0
	}
	return factorialHelper(n, 1)
}

func factorialHelper(n, result int) int {
	if n == 0 {
		return result
	}
	return factorialHelper(n-1, n*result)
}

func BenchmarkFactorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		factorial(12)
	}
}
