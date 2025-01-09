package main

import (
	"runtime"
	"sync"
	"testing"
)

func sumAndDiffConcurrently(a, b int) (int, int) {
	sum := make(chan int)
	diff := make(chan int)

	go func() {
		sum <- a + b
	}()

	go func() {
		diff <- a - b
	}()

	return <-sum, <-diff
}

func TestSumAndDiffConcurrently(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Utilize all available CPU cores

	testing.Benchmark(func(b *testing.B) {
		var wg sync.WaitGroup
		wg.Add(b.N)

		for i := 0; i < b.N; i++ {
			go func(i int) {
				_, _ = sumAndDiffConcurrently(i, b.N-i)
				wg.Done()
			}(i)
		}

		wg.Wait()
	})
}
