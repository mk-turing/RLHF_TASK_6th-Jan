package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu       sync.RWMutex
	data     map[string]interface{}
	expiry   map[string]time.Time
	ttl      time.Duration
	quit     chan struct{}
	wg       sync.WaitGroup
	expiryMu sync.Mutex // New lock for expiry operations
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		data:   make(map[string]interface{}),
		expiry: make(map[string]time.Time),
		ttl:    ttl,
		quit:   make(chan struct{}),
	}
	c.wg.Add(1)
	go c.expire()
	return c
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expiry[key] = time.Now().Add(c.ttl)
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.data[key]
	if !ok {
		return nil, false
	}

	expiry, exists := c.expiry[key]
	if !exists || time.Now().After(expiry) {
		return nil, false
	}

	return v, true
}

func (c *Cache) expire() {
	defer c.wg.Done()
	for {
		select {
		case <-time.After(c.ttl / 2):
			c.expiryMu.Lock() // Acquire the new lock here
			now := time.Now()
			for key, exp := range c.expiry {
				if exp.Before(now) {
					c.mu.Lock() // Lock the data structure for modification
					delete(c.data, key)
					delete(c.expiry, key)
					c.mu.Unlock()
				}
			}
			c.expiryMu.Unlock() // Release the lock after expiration
		case <-c.quit:
			return
		}
	}
}

func (c *Cache) Close() {
	c.quit <- struct{}{}
	c.wg.Wait()
}

func main() {
	cache := NewCache(time.Second * 2)
	defer cache.Close()

	// Set an item in the cache with a TTL
	cache.Set("mykey", "Hello World!")

	// Get the item from the cache
	value, found := cache.Get("mykey")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}

	time.Sleep(time.Second * 3)

	// The item should have expired
	value, found = cache.Get("mykey")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}
}
