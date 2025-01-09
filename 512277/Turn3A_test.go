package _12277

import (
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkServerPerformance(b *testing.B) {
	// Get the number of workers from the environment variable
	numWorkersStr := os.Getenv("NUM_WORKERS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		b.Fatalf("Invalid NUM_WORKERS environment variable: %v", err)
	}

	// Get the traffic pattern from the environment variable
	trafficPatternStr := os.Getenv("TRAFFIC_PATTERN")
	trafficPattern, err := strconv.Atoi(trafficPatternStr)
	if err != nil {
		b.Fatalf("Invalid TRAFFIC_PATTERN environment variable: %v", err)
	}

	// Get the CPU limit (percentage) from the environment variable
	cpuLimitStr := os.Getenv("CPU_LIMIT")
	cpuLimit, err := strconv.ParseFloat(cpuLimitStr, 64)
	if err != nil {
		b.Fatalf("Invalid CPU_LIMIT environment variable: %v", err)
	}

	// Get the memory limit (MB) from the environment variable
	memLimitStr := os.Getenv("MEM_LIMIT")
	memLimit, err := strconv.ParseInt(memLimitStr, 10, 64)
	if err != nil {
		b.Fatalf("Invalid MEM_LIMIT environment variable: %v", err)
	}

	// Implement CPU throttling (simulated)
	runtime.GOMAXPROCS(int(cpuLimit / 100 * float64(runtime.NumCPU())))

	// Implement memory cap (simulated)
	if memLimit > 0 {
		log.Printf("Simulating memory limit of %d MB\n", memLimit)
		// Simulate high memory usage
		maxMemory := 100 * 1024 * 1024
		allocData := make([]byte, maxMemory)
		copy(allocData, allocData)
	}

	// Simulate network latency
	networkLatencyStr := os.Getenv("NETWORK_LATENCY")
	networkLatency, err := strconv.Atoi(networkLatencyStr)
	if err != nil {
		networkLatency = 0
	}
	networkLatencyInMS := time.Duration(networkLatency) * time.Millisecond

	// Create worker pool and load balancer
	var wg sync.WaitGroup
	workers := make([]chan struct{}, numWorkers)
	for i := range workers {
		workers[i] = make(chan struct{}, 1) // Buffered channel with size 1
		wg.Add(1)
		go func(id int) {
			defer wg.Done() // Ensure Done is called
			for {
				select {
				case <-workers[id]:
					time.Sleep(time.Millisecond * 10) // Simulate work duration
					workers[id] <- struct{}{}         // Signal that work is complete
				}
			}
		}(i)
	}

	// Start the benchmark timer
	b.StartTimer()

	// Simulate load balancing and traffic patterns
	for n := 0; n < b.N; n++ {
		time.Sleep(networkLatencyInMS) // Simulate network latency

		// Round-robin load balancing
		currentWorker := n % numWorkers
		workers[currentWorker] <- struct{}{} // Assign work to the worker

		// Simulate client request based on traffic pattern
		switch trafficPattern {
		case 1: // Uniform traffic
		case 2: // Bursty traffic (every 5 requests)
			if n%5 == 0 {
				// Simulate burst of 5 requests to the current worker
				for i := 0; i < 5; i++ {
					workers[currentWorker] <- struct{}{}
				}
			}
		default:
			log.Println("Invalid traffic pattern")
		}

		// Release worker after simulating work
		workers[currentWorker] <- struct{}{} // Ensure worker gets to work again
	}

	// Stop the benchmark timer
	b.StopTimer()

	// Wait for all workers to finish
	wg.Wait()

	// Log performance metrics at the end of the benchmark
	log.Printf("Benchmark completed: %s", b.Name())

	// Use b.Elapsed() to get the total benchmark duration
	duration := b.Elapsed()

	// Calculate operations per second and average time per operation
	if duration > 0 {
		opsPerSecond := float64(b.N) / duration.Seconds()
		avgTimePerOp := duration.Nanoseconds() / int64(b.N)
		log.Printf("Operations per second: %.2f", opsPerSecond)
		log.Printf("Average time per operation: %d ns", avgTimePerOp)
	} else {
		log.Println("Benchmark duration is 0, unable to calculate operations per second or average time per operation.")
	}
}
