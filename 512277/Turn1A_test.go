package _12277

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkServerPerformance(b *testing.B) {
	// Get environment variables
	concurrencyStr := os.Getenv("CONCURRENCY")
	requestSizeStr := os.Getenv("REQUEST_SIZE")
	intervalStr := os.Getenv("INTERVAL")

	var concurrency int
	var requestSize int
	var interval time.Duration

	// Convert environment variables to usable types
	if concurrencyStr != "" {
		concurrency, _ = strconv.Atoi(concurrencyStr)
	}

	if requestSizeStr != "" {
		requestSize, _ = strconv.Atoi(requestSizeStr)
	}

	if intervalStr != "" {
		interval, _ = time.ParseDuration(intervalStr)
	}

	// Simulate load
	simulateLoad(concurrency, requestSize, interval)
}

func simulateLoad(concurrency int, requestSize int, interval time.Duration) {
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				time.Sleep(interval)
				// Simulate request
				// You could send a request to a server or perform a CPU-intensive task
				fmt.Println("Sending request with size", requestSize)
			}
		}()
	}

	wg.Wait()
}
