package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Your function to be benchmarked (modified to be concurrency-safe)
func myFunction(input int, wg *sync.WaitGroup) int {
	defer wg.Done()
	// Function logic that can handle multiple concurrent requests
	// For example, a simple calculation with a sleep to simulate computation
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	return input * input
}

// BenchmarkMyFunction_Concurrency tests the system's concurrency handle and thread safety
func BenchmarkMyFunction_Concurrency(b *testing.B) {
	const numWorkers = 100
	const maxInput = 1000000

	for numRequests := 100; numRequests <= 10000; numRequests *= 10 {
		b.Run(fmt.Sprintf("NumRequests_%d", numRequests), func(b *testing.B) {
			c := concurrency.New(numWorkers, numRequests)
			c.Run(func() error {
				wg := sync.WaitGroup{}
				for i := 0; i < numRequests; i++ {
					wg.Add(1)
					go func() {
						input := rand.Intn(maxInput)
						myFunction(input, &wg)
					}()
				}
				wg.Wait()
				return nil
			})
			b.ReportAllocs()
			b.ReportMetric(float64(numRequests)/c.Duration().Seconds(), "requests/s")
		})
	}
}

// Benchmarks under the concurrent flow will simulate
// multiple concurrent requests and test for any race conditions.

func main() {
	testing.Main()
}
