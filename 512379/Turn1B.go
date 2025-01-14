package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	data map[string]string
	mux  sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}
func (c *Cache) Set(key, value string, expiration time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.data[key] = value
	go func() {
		<-time.After(expiration)
		c.Delete(key)
	}()
}
func (c *Cache) Get(key string) (string, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	value, ok := c.data[key]
	return value, ok
}
func (c *Cache) Delete(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.data, key)
}

// Simple microservice using the cache
type MicroService struct {
	cache *Cache
}

func NewMicroService(cache *Cache) *MicroService {
	return &MicroService{cache: cache}
}
func (ms *MicroService) GetData(key string) (string, bool) {
	value, ok := ms.cache.Get(key)
	if ok {
		fmt.Println("Retrieved data from cache:", key, value)
		return value, true
	}
	fmt.Println("Data not found in cache, fetching from external service...")
	// Simulate fetching data from an external service
	time.Sleep(2 * time.Second)
	value = "External Service Data"
	ms.cache.Set(key, value, 5*time.Second) // Set expiration to 5 seconds
	return value, true
}

func main() {
	// Create a new cache instance
	cache := NewCache()

	// Create microservices with the cache
	ms1 := NewMicroService(cache)
	ms2 := NewMicroService(cache)

	// Request data from microservices, they will use the cache internally
	key := "user_data"
	data, _ := ms1.GetData(key)
	data, _ = ms2.GetData(key) // Data will be retrieved from the cache in the second request
}
