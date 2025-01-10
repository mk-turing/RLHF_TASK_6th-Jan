package _12283

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Simulate variable resource allocation within the worker goroutines
func doSomeWorkWithResourceAllocation(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond) // Simulate variable work
	// Perform resource-allocating tasks like disk reads/writes or memory allocations
}

// Benchmark tests with variable resource allocation and multi-threading
func BenchmarkVariableResourceAllocation(b *testing.B) {
	const numWorkers = 10
	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(numWorkers)
		for j := 0; j < numWorkers; j++ {
			go doSomeWorkWithResourceAllocation(&wg)
		}
		wg.Wait()
	}
}

// Simulate random network latency and occasional network failures in a more complex way
func BenchmarkMoreComplexNetworkScenario(b *testing.B) {
	const numWorkers = 5
	var wg sync.WaitGroup

	// Simulate various network conditions by creating a network latency channel
	networkLatencyChannel := make(chan time.Duration, numWorkers)

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
			networkLatencyChannel <- time.Duration(rand.Intn(100)) * time.Millisecond
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(numWorkers)
		for j := 0; j < numWorkers; j++ {
			go func(networkLatencyChannel chan time.Duration) {
				defer wg.Done()
				// Simulate network latency
				latency := <-networkLatencyChannel
				time.Sleep(latency)
				doSomeWork()
			}(networkLatencyChannel)
		}
		wg.Wait()
	}
}
