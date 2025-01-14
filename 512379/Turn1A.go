package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu     sync.RWMutex // For thread safety
	data   map[string]interface{}
	expiry map[string]time.Time // Store expiration times separately
	ttl    time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		data:   make(map[string]interface{}),
		expiry: make(map[string]time.Time),
		ttl:    ttl,
	}
}

func (c *Cache) Set(key string, value interface{}, expiresAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
	c.expiry[key] = expiresAt
	go c.expire(key, expiresAt)
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

func (c *Cache) expire(key string, expiresAt time.Time) {
	time.Sleep(time.Until(expiresAt))
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	delete(c.expiry, key)
}

func main() {
	cache := NewCache(time.Second * 2)

	// Set an item in the cache with a TTL
	expiresAt := time.Now().Add(time.Second * 1)
	cache.Set("mykey", "Hello World!", expiresAt)

	// Get the item from the cache
	value, found := cache.Get("mykey")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}

	time.Sleep(time.Second * 2)

	// The item should have expired
	value, found = cache.Get("mykey")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}
}
