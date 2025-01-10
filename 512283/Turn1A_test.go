package _12283

import (
	"testing"
	"time"
)

// Example function to simulate some work
func doSomeWork() {
	time.Sleep(1 * time.Millisecond) // Simulate some work taking time
}

// Benchmark function that calls doSomeWork
func BenchmarkDoSomeWork(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doSomeWork()
	}
}

// Example function to simulate another type of work
func doMoreWork() {
	time.Sleep(2 * time.Millisecond) // Simulate more work taking more time
}

// Benchmark function that calls doMoreWork
func BenchmarkDoMoreWork(b *testing.B) {
	for i := 0; i < b.N; i++ {
		doMoreWork()
	}
}
