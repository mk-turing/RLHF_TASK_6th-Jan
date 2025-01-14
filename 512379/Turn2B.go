package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu     sync.RWMutex
	data   map[string]interface{}
	expiry map[string]time.Time
	ttl    time.Duration
	stop   chan struct{}
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		data:   make(map[string]interface{}),
		expiry: make(map[string]time.Time),
		ttl:    ttl,
		stop:   make(chan struct{}),
	}
	// Start the goroutine to handle expiration
	go c.expirationWorker()
	return c
}
func (c *Cache) expirationWorker() {
	ticker := time.NewTicker(c.ttl / 10) // Check expiration every 10th of TTL
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, expiresAt := range c.expiry {
				if now.After(expiresAt) {
					delete(c.data, key)
					delete(c.expiry, key)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			ticker.Stop()
			return
		}
	}
}
func (c *Cache) Set(key string, value interface{}, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expiry[key] = expiresAt
}
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.data[key]
	if !ok {
		return nil, false
	}

	expiresAt, exists := c.expiry[key]
	if !exists || time.Now().After(expiresAt) {
		return nil, false
	}

	return v, true
}
func (c *Cache) Close() {
	close(c.stop)
}

func main() {
	cache := NewCache(time.Second * 2)

	// Set an item in the cache with a TTL
	expiresAt := time.Now().Add(time.Second * 1)
	cache.Set("mykey", "Hello World!", expiresAt)
	time.Sleep(time.Second * 1)
	value, found := cache.Get("mykey")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}
	// Wait for the entry to expire and then get it again
	time.Sleep(time.Second * 2)
	value, found = cache.Get("mykey")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss (expired).")
	}
	cache.Close()
}
