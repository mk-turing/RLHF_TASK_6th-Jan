package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type ValidationCallback func(data interface{}) error

type RateLimiter struct {
	burst     int
	interval  time.Duration
	lastCall  time.Time
	callCount int
	lock      sync.Mutex
}

func NewRateLimiter(burst int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		burst:    burst,
		interval: interval,
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	now := time.Now()
	if rl.callCount < rl.burst || now.Sub(rl.lastCall) > rl.interval {
		rl.lastCall = now
		rl.callCount = 1
		return true
	}

	rl.callCount++
	return false
}

func ValidateData(data interface{}, cb ValidationCallback) error {
	// Your actual validation logic here
	time.Sleep(time.Millisecond * 50)

	return cb(data)
}

func worker(id int, data []interface{}, rateLimiter *RateLimiter, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, item := range data {
		if rateLimiter.Allow() {
			err := ValidateData(item, func(data interface{}) error {
				fmt.Printf("Worker %d validated data: %v\n", id, data)
				return nil
			})
			if err != nil {
				fmt.Printf("Worker %d encountered error: %v\n", id, err)
			}
		} else {
			fmt.Printf("Worker %d: Rate limit exceeded, skipped validation for %v\n", id, item)
		}
	}
}

func adjustRateLimits(rateLimiter *RateLimiter) {
	ticker := time.NewTicker(time.Second * 5) // Adjust every 5 seconds
	defer ticker.Stop()

	for range ticker.C {
		loadAvg := runtime.LoadAvg()
		fmt.Printf("Current system load average: %v\n", loadAvg)

		// Adjust the burst rate based on load average
		// Example logic: Increase burst rate if load is low, decrease if load is high
		if loadAvg.Avg1 < 1.5 {
			rateLimiter.burst = 5 // High burst rate for low load
		} else if loadAvg.Avg1 < 3.0 {
			rateLimiter.burst = 3 // Moderate burst rate for medium load
		} else {
			rateLimiter.burst = 1 // Low burst rate for high load
		}

		fmt.Printf("Adjusted rate limit: burst=%d\n", rateLimiter.burst)
	}
}

func main() {
	// Sample data
	data := []interface{}{"data1", "data2", "data3", "data4", "data5", "data6", "data7", "data8", "data9", "data10"}

	numWorkers := 3
	batchSize := len(data) / numWorkers

	rateLimiter := NewRateLimiter(2, time.Second) // Initial rate limit

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Start a goroutine to adjust rate limits periodically
	go adjustRateLimits(rateLimiter)

	for i := 0; i < numWorkers; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(data) {
			end = len(data)
		}
		go worker(i+1, data[start:end], rateLimiter, &wg)
	}

	wg.Wait()
	fmt.Println("All data validation completed.")
}
