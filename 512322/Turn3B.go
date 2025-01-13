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

type RetryStrategy struct {
	MaxRetries int
	Backoff    time.Duration
	Jitter     float64
}

func (rs RetryStrategy) BackoffDuration(attempt int) time.Duration {
	d := rs.Backoff * time.Duration(2<<uint(attempt-1))
	jitter := time.Duration(float64(d) * rs.Jitter * rand.Float64())
	return d + jitter
}

func validateDataWithRetry(data interface{}, cb ValidationCallback, rl *RateLimiter, rs RetryStrategy) error {
	for attempt := 1; attempt <= rs.MaxRetries; attempt++ {
		if rl.Allow() {
			err := cb(data)
			if err == nil {
				return nil // Success, no need to retry
			}
			fmt.Printf("Validation failed for data: %v, attempt %d/%d, retrying in %s\n", data, attempt, rs.MaxRetries, rs.BackoffDuration(attempt))
		} else {
			fmt.Printf("Rate limit exceeded, skipping validation for %v and retrying later\n", data)
		}
		time.Sleep(rs.BackoffDuration(attempt))
	}
	return fmt.Errorf("exhausted retries for data: %v", data)
}

func worker(id int, data []interface{}, rateLimiter *RateLimiter, retryStrategy RetryStrategy, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var lastAlloc uint64
	for _, item := range data {
		if err := validateDataWithRetry(item, func(data interface{}) error {
			// Your actual validation logic here
			time.Sleep(time.Millisecond * 50)
			return nil
		}, rateLimiter, retryStrategy); err != nil {
			fmt.Printf("Worker %d: Failed to validate data %v after %d retries\n", id, item, retryStrategy.MaxRetries)
		}