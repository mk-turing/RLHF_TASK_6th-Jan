package main

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Your function to be benchmarked
func myFunction(input int) int {
	// Function logic that can handle extreme inputs
	// For example, a simple calculation with a sleep to simulate computation
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	return input * input
}

// BenchmarkMyFunction_Normal loads the system with normal requests
func BenchmarkMyFunction_Normal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		myFunction(i)
	}
}

// BenchmarkMyFunction_Stress tests the system with extreme inputs and high load
func BenchmarkMyFunction_Stress(b *testing.B) {
	const numWorkers = 100
	const maxInput = 1000000

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < b.N; j++ {
				// Generate random extreme inputs
				input := rand.Intn(maxInput)
				myFunction(input)
			}
		}()
	}

	wg.Wait()
}

// BenchmarkMyFunction_Extreme tests the system with the maximum possible input
func BenchmarkMyFunction_Extreme(b *testing.B) {
	const maxInput = 1000000
	for i := 0; i < b.N; i++ {
		myFunction(maxInput)
	}
}

// BenchmarkMyFunction_Small tests the system with the minimum possible input
func BenchmarkMyFunction_Small(b *testing.B) {
	for i := 0; i < b.N; i++ {
		myFunction(1)
	}
}
