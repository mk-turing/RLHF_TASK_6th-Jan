package _93902

import (
	"testing"
)

// memoizedFactorial stores previously calculated factorials
var memoizedFactorial = make(map[int]int)

// factorial calculates the factorial of a number n using memoization
func factorial(n int) int {
	// Base case: The factorial of 0 and 1 is 1.
	if n < 2 {
		return 1
	}
	// Check if the factorial of n has already been calculated and stored.
	if result, ok := memoizedFactorial[n]; ok {
		return result
	}
	// Recursive case: Calculate the factorial of n and store it in the memoization map.
	result := n * factorial(n-1)
	memoizedFactorial[n] = result
	return result
}

// BenchmarkFactorial benchmarks the factorial function with memoization.
func BenchmarkFactorial(b *testing.B) {
	// Clear the memoizedFactorial map before each benchmark run.
	memoizedFactorial = make(map[int]int)
	for i := 0; i < b.N; i++ {
		factorial(15)
	}
}

func TestFactorial(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Small Factorial", args{3}, 6},
		{"Large Factorial", args{12}, 479001600},
		{"Negative Factorial", args{-5}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := factorial(tt.args.n); got != tt.want {
				t.Errorf("factorial() = %v, want %v", got, tt.want)
			}
		})
	}
}
