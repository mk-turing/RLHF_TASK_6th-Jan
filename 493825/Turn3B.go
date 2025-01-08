package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func sumRangeManual(start, end int) int {
	total := 0
	for i := start; i <= end; i++ {
		total += i
	}
	return total
}

func parallelSumRange(start, end int, wg *sync.WaitGroup) {
	defer wg.Done()
	total := sumRangeManual(start, end)
	fmt.Printf("Partial Sum (%d-%d): %d\n", start, end, total)
}

func benchmarkParallelSumRange(threads int) {
	var wg sync.WaitGroup

	start := time.Now()

	rangeSize := 1000000
	partitionSize := rangeSize / threads

	wg.Add(threads)
	for i := 0; i < threads; i++ {
		startPartition := i*partitionSize + 1
		endPartition := min((i+1)*partitionSize, rangeSize)
		go parallelSumRange(startPartition, endPartition, &wg)
	}
	wg.Wait()

	end := time.Now()
	fmt.Printf("Parallel (%d threads) Benchmark: %.6f seconds\n", threads, end.Sub(start).Seconds())
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	fmt.Println("Single-Threaded Benchmark:")
	benchmarkSingleThreaded()

	fmt.Println("\nParallel Benchmark (2 threads):")
	benchmarkParallelSumRange(2)

	fmt.Println("\nParallel Benchmark (4 threads):")
	benchmarkParallelSumRange(4)

	fmt.Println("\nParallel Benchmark (8 threads):")
	benchmarkParallelSumRange(8)
}

func benchmarkSingleThreaded() {
	runtime.GOMAXPROCS(1)
	start := time.Now()
	sumRangeManual(1, 1000000)
	end := time.Now()
	fmt.Printf("Single Thread Benchmark: %.6f seconds\n", end.Sub(start).Seconds())
}
