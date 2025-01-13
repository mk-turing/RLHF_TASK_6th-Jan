package main

import (
	"container/heap"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type ValidationCallback func(data interface{}) error

type Callback struct {
	Priority int
	Data     interface{}
	Callback ValidationCallback
}

// PriorityQueue implements heap.Interface and holds Callbacks with priorities.
type PriorityQueue []*Callback

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Higher priority callbacks have lower indices
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Callback))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

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

func worker(id int, pq *PriorityQueue, rateLimiter *RateLimiter, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var lastAlloc uint64

	for {
		if pq.Len() == 0 {
			// If the queue is empty, break the loop
			break
		}

		// Take the highest priority callback from the queue
		c := heap.Pop(pq).(*Callback)
		if c == nil || c.Callback == nil {
			// If c or c.Callback is nil, continue to prevent dereferencing a nil pointer
			continue
		}

		err := RetryCallback(rateLimiter, c.Callback, c.Data, 3) // Max of 3 retries
		if err != nil {
			fmt.Printf("Worker %d: Validation failed for %v: %v\n", id, c.Data, err)
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
	// Sample data with varying priorities
	data := []*Callback{
		{Priority: 1, Data: "data1"},
		{Priority: 2, Data: "data2"},
		{Priority: 3, Data: "data3"},
		{Priority: 2, Data: "data4"},
		{Priority: 1, Data: "data5"},
		{Priority: 3, Data: "data6"},
		{Priority: 2, Data: "data7"},
		{Priority: 1, Data: "data8"},
		{Priority: 3, Data: "data9"},
		{Priority: 2, Data: "data10"},
	}

	numWorkers := 3
	batchSize := len(data) / numWorkers

	rateLimiter := NewRateLimiter(10, time.Second) // Allow 10 calls per second initially
	pq := &PriorityQueue{}
	heap.Init(pq)

	for _, item := range data {
		heap.Push(pq, item)
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		end := (i + 1) * batchSize
		if end > len(data) {
			end = len(data)
		}
		go worker(i+1, pq, rateLimiter, &wg)
	}
	wg.Wait()
	fmt.Println("All data validation completed.")
}
