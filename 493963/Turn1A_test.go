package main

import (
	"testing"
)

func BenchmarkMyFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Your function to be benchmarked
		myFunction()
	}
}

func myFunction() {
	// Function logic
}
