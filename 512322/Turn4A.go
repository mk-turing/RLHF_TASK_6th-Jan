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

func RetryCallback(rateLimiter *RateLimiter, cb ValidationCallback, data interface{}, maxAttempts int, priority int) error {
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d of %d for data with priority %d: %v\n", attempt, maxAttempts, priority, data)

		if rateLimiter.Allow() {
			err := cb(data)
			if err == nil {
				return nil
			}
			fmt.Printf("Callback with priority %d failed, error: %v\n", priority, err)
		} else {
			fmt.Printf("Rate limit exceeded, skipping retry for data with priority %d: %v\n", priority, data)
			return fmt.Errorf("Rate limit exceeded")
		}

		backoffTime := time.Duration(1<<(attempt-1)) * time.Second // Exponential backoff
		time.Sleep(backoffTime)
	}

	return fmt.Errorf("Max attempts reached, validation failed for %v with priority %d", data, priority)
}

// Implement a priority queue using a slice with atomic swaps for fairness
type PriorityQueue struct {
	queue []*callbackItem
	mutex sync.RWMutex
}

type callbackItem struct {
	priority int
	callback func(interface{})
	data     interface{}
}

func (pq *PriorityQueue) Add(priority int, callback ValidationCallback, data interface{}) {
	pq.mutex.Lock()
	defer pq.mutex.Unlock()

	item := &callbackItem{
		priority: priority,
		callback: callback,
		data:     data,
	}
	// Insert item at the correct position to maintain priority order
	i := 0
	for i < len(pq.queue) && pq.queue[i].priority < item.priority {
		i++
	}
	pq.queue = append(pq.queue[:i], item, pq.queue[i:]...)
}

func (pq *PriorityQueue) Pop() (int, func(interface{}) error, interface{}) {
	pq.mutex.RLock()
	defer pq.mutex.RUnlock()

	if len(pq.queue) == 0 {
		return 0, nil, nil
	}
	return pq.queue[0].priority, pq.queue[0].callback, pq.queue[0].data
}

func worker(id int, dataQueue *PriorityQueue, rateLimiter *RateLimiter, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastAlloc uint64
	for {
		priority, callback, data := dataQueue.Pop()
		if priority == 0 {
			break // Stop the worker if there are no more items in the queue
		}

		err := RetryCallback(rateLimiter, callback, data, 3, priority) // Max of 3 retries
		if err != nil {
			fmt.Printf("Worker %d: Validation failed for data with priority %d: %v: %v\n", id, priority, data, err)
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

	dataQueue := &PriorityQueue{}
	rateLimiter := NewRateLimiter(10, time.Second) // Allow 10 calls per second initially

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Add data to the queue with different priorities
	for i, item := range data {
		if i < len(data)/2 {
			dataQueue.Add(1, func(data interface{}) error {
				return ValidateData(data)
			}, item) // Higher priority for the first half of the data
		} else {
			dataQueue.Add(2, func(data interface{}) error {
				return ValidateData(data)
			}, item) // Lower priority for the second half of the data
		}
	}

	for i := 0; i < numWorkers; i++ {
		go worker(i+1, dataQueue, rateLimiter, &wg)
	}

	wg.Wait()
	fmt.Println("All data validation completed.")
}
