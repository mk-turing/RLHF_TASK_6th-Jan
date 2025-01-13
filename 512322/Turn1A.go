package main

import (
	"sync"
	"time"
)

// Callback represents a single validation callback function
type Callback func()

// RateLimiter controls the rate of callback executions
type RateLimiter struct {
	// Limit specifies the maximum number of callbacks to allow per interval
	Limit int
	// Interval is the time frame over which the limit is applied
	Interval time.Duration
	// Checkpoint records the last time a callback was executed
	Checkpoint time.Time
	// MaxCalls stores the number of callbacks executed since the last checkpoint reset
	MaxCalls int
	// Lock ensures thread-safe access to the RateLimiter's state
	sync.Mutex
}

// CanExecute determines if a callback can be executed within the rate limit
func (rl *RateLimiter) CanExecute() bool {
	rl.Lock()
	defer rl.Unlock()

	now := time.Now()
	if now.Sub(rl.Checkpoint) >= rl.Interval {
		rl.MaxCalls = 0
		rl.Checkpoint = now
	}

	if rl.MaxCalls >= rl.Limit {
		return false
	}

	rl.MaxCalls++
	return true
}

// CallbackManager handles the registration and execution of callbacks with rate limiting
type CallbackManager struct {
	// Limiter provides the rate-limiting mechanism
	Limiter *RateLimiter
	// Callbacks is a list of registered callbacks
	Callbacks []Callback
	// Queue ensures that callbacks are executed in order
	Queue <-chan Callback
}

// AddCallback registers a new callback function
func (cm *CallbackManager) AddCallback(cb Callback) {
	cm.Callbacks = append(cm.Callbacks, cb)
}

// ExecuteCallbacks starts the loop to process callbacks while respecting the rate limit
func (cm *CallbackManager) ExecuteCallbacks() {
	for cb := range cm.Queue {
		if cm.Limiter != nil && !cm.Limiter.CanExecute() {
			// Slow down execution if the limit is exceeded
			time.Sleep(1 * time.Second)
		}

		cb()
	}
}

func main() {
	limiter := &RateLimiter{Limit: 10, Interval: 1 * time.Second}
	cm := &CallbackManager{Limiter: limiter, Queue: make(chan Callback)}

	// Register some validation callbacks
	cm.AddCallback(func() {
		println("Validating data entry...")
		time.Sleep(100 * time.Millisecond) // Simulate processing time
	})
	cm.AddCallback(func() {
		println("Validating user credentials...")
		time.Sleep(200 * time.Millisecond) // Simulate processing time
	})
	cm.AddCallback(func() {
		println("Validating transaction details...")
		time.Sleep(150 * time.Millisecond) // Simulate processing time
	})

	// Start processing callbacks
	go cm.ExecuteCallbacks()

	// Send registered callbacks to the queue
	for _, cb := range cm.Callbacks {
		cm.Queue <- cb
	}

	// Close the queue after all callbacks are added
	close(cm.Queue)
}
