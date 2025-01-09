package _12291

import (
	"testing"
)

// Function with no arguments and no return value
func noArgNoReturn() {
}

// Function with one argument and no return value
func oneArgNoReturn(arg int) {
}

// Function with no arguments and one return value
func noArgOneReturn() int {
	return 0
}

// Function with one argument and one return value
func oneArgOneReturn(arg int) int {
	return arg
}

// Function with two arguments and no return value
func twoArgsNoReturn(arg1, arg2 int) {
}

// Function with two arguments and one return value
func twoArgsOneReturn(arg1, arg2 int) int {
	return arg1 + arg2
}

// Function with variadic arguments and no return value
func variadicArgsNoReturn(args ...int) {
}

// Function with variadic arguments and one return value
func variadicArgsOneReturn(args ...int) int {
	sum := 0
	for _, arg := range args {
		sum += arg
	}
	return sum
}

// Benchmark function for no arguments and no return value
func BenchmarkNoArgNoReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		noArgNoReturn()
	}
}

// Benchmark function for one argument and no return value
func BenchmarkOneArgNoReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		oneArgNoReturn(1)
	}
}

// Benchmark function for no arguments and one return value
func BenchmarkNoArgOneReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		noArgOneReturn()
	}
}

// Benchmark function for one argument and one return value
func BenchmarkOneArgOneReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		oneArgOneReturn(1)
	}
}

// Benchmark function for two arguments and no return value
func BenchmarkTwoArgsNoReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		twoArgsNoReturn(1, 2)
	}
}

// Benchmark function for two arguments and one return value
func BenchmarkTwoArgsOneReturn(b *testing.B) {
	for i := 0; i < b.N; i++ {
		twoArgsOneReturn(1, 2)
	}
}

// Benchmark function for variadic arguments and no return value
func BenchmarkVariadicArgsNoReturn(b *testing.B) {
	args := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		variadicArgsNoReturn(args...)
	}
}

// Benchmark function for variadic arguments and one return value
func BenchmarkVariadicArgsOneReturn(b *testing.B) {
	args := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		variadicArgsOneReturn(args...)
	}
}

func main() {
	// Run all the benchmarks
	testing.Main(func(m *testing.M) {
		return m.Run()
	}, []testing.InternalTest{
		{Name: "BenchmarkNoArgNoReturn"},
		{Name: "BenchmarkOneArgNoReturn"},
		{Name: "BenchmarkNoArgOneReturn"},
		{Name: "BenchmarkOneArgOneReturn"},
		{Name: "BenchmarkTwoArgsNoReturn"},
		{Name: "BenchmarkTwoArgsOneReturn"},
		{Name: "BenchmarkVariadicArgsNoReturn"},
		{Name: "BenchmarkVariadicArgsOneReturn"},
	}, nil, nil)
}
