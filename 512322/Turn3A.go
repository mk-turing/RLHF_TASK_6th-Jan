package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type ValidationCallback func(data interface{}) error

type RateLimiter struct {
	burst        int
	interval     time.Duration
	lastCall     time.Time
	callCount    int
	currentLimit int
	lock         sync.Mutex
}

func NewRateLimiter(burst int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		burst:    burst,
		interval: interval,
	}
	rl.currentLimit = burst
	return rl
}

func (rl *RateLimiter) Allow() bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	now := time.Now()
	if rl.callCount < rl.currentLimit || now.Sub(rl.lastCall) > rl.interval {
		rl.lastCall = now
		rl.callCount = 1
		return true
	}

	rl.callCount++
	return false
}

func (rl *RateLimiter) AdjustRateLimit(load float64) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	newLimit := int(float64(rl.burst) * load)
	if newLimit < 1 {
		newLimit = 1
	}
	rl.currentLimit = newLimit
	fmt.Printf("Adjusting rate limit to: %d\n", newLimit)
}

func measureLoad(lastAlloc uint64) float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate the difference in allocated memory since the last check
	allocRate := float64(m.Alloc-lastAlloc) / float64(1024*1024) // MB per second

	return allocRate
}

func ValidateData(data interface{}) error {
	// Your actual validation logic here
	// For this example, we introduce random failures to simulate real-world scenarios
	if rand.Intn(2) == 0 {
		fmt.Printf("Validation failed for: %v\n", data)
		return fmt.Errorf("Validation failed for %v", data)
	}

	time.Sleep(time.Millisecond * 50)
	fmt.Printf("Data validated successfully: %v\n", data)
	return nil
}

func RetryCallback(rateLimiter *RateLimiter, cb ValidationCallback, data interface{}, maxAttempts int) error {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d of %d for data: %v\n", attempt, maxAttempts, data)

		if rateLimiter.Allow() {
			err := cb(data)
			if err == nil {
				return nil
			}
			fmt.Printf("Callback failed, error: %v\n", err)
		} else {
			fmt.Printf("Rate limit exceeded, skipping retry for data: %v\n", data)
			return fmt.Errorf("Rate limit exceeded")
		}

		backoffTime := time.Duration(1<<(attempt-1)) * time.Second // Exponential backoff
		time.Sleep(backoffTime)
	}

	return fmt.Errorf("Max attempts reached, validation failed for %v", data)
}

func worker(id int, data []interface{}, rateLimiter *RateLimiter, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastAlloc uint64
	for _, item := range data {
		err := RetryCallback(rateLimiter, func(data interface{}) error {
			return ValidateData(data)
		}, item, 3) // Max of 3 retries
		if err != nil {
			fmt.Printf("Worker %d: Validation failed for %v: %v\n", id, item, err)
		}

		select {
		case <-ticker.C:
			// Measure memory allocation rate and adjust the rate limit
			load := measureLoad(lastAlloc)
			rateLimiter.AdjustRateLimit(load)
			lastAlloc = runtime.MemStats{}.Alloc // Update the lastAlloc with the current memory usage
		default:
		}
	}
}

func main() {
	// Sample data
	data := []interface{}{"data1", "data2", "data3", "data4", "data5", "data6", "data7", "data8", "data9", "data10"}

	numWorkers := 3
	batchSize := len(data) / numWorkers

	rateLimiter := NewRateLimiter(10, time.Second) // Allow 10 calls per second initially

	var wg sync.WaitGroup
	wg.Add(numWorkers)

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
