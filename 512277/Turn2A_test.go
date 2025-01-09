package _12277

import (
	"log"
	"os"
	"runtime"
	"strconv"
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

	// Create a worker pool with the specified number of workers
	workers := make(chan struct{}, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers <- struct{}{}
	}

	// Start the benchmark timer
	b.StartTimer()

	// Simulate work using a goroutine
	go func() {
		for {
			<-workers
			time.Sleep(time.Millisecond * 10) // Simulate work duration
			workers <- struct{}{}
		}
	}()

	// Simulate network latency before each request
	for n := 0; n < b.N; n++ {
		time.Sleep(networkLatencyInMS) // Simulate network latency
		// Simulate client requests by acquiring a worker
		<-workers
		// Release the worker after simulating work
		workers <- struct{}{}
	}

	// Stop the benchmark timer
	b.StopTimer()

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
