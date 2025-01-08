package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
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

func benchmarkParallelSumRange(threads int, wg *sync.WaitGroup) {
	start := time.Now()

	rangeSize := 1000000
	partitionSize := rangeSize / threads

	for i := 0; i < threads; i++ {
		startPartition := i*partitionSize + 1
		endPartition := min((i+1)*partitionSize, rangeSize)
		go parallelSumRange(startPartition, endPartition, &wg)
	}
	wg.Wait()

	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Parallel (%d threads) Benchmark: %.6f seconds\n", threads, duration.Seconds())
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
	var wg sync.WaitGroup
	wg.Add(1)
	go benchmarkParallelSumRange(2, &wg)
	wg.Wait()

	fmt.Println("\nParallel Benchmark (4 threads):")
	wg.Add(1)
	go benchmarkParallelSumRange(4, &wg)
	wg.Wait()

	fmt.Println("\nParallel Benchmark (8 threads):")
	wg.Add(1)
	go benchmarkParallelSumRange(8, &wg)
	wg.Wait()

	runProfiles()
}

func benchmarkSingleThreaded() {
	runtime.GOMAXPROCS(1)
	start := time.Now()
	sumRangeManual(1, 1000000)
	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("Single Thread Benchmark: %.6f seconds\n", duration.Seconds())
}

// Start a server to export the profile
func runProfiles() {
	// Start a server for the CPU profile.
	go http.HandleFunc("/pprof/", pprof.Index)
	go http.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	go http.HandleFunc("/pprof/profile", pprof.Profile)
	go http.HandleFunc("/pprof/symbol", pprof.Symbol)

	// CPU Profile
	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()

	fmt.Println("Profiling server running on :8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
