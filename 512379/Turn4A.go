package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

// Node represents a node in the doubly linked list
type Node struct {
	key        string
	value      interface{}
	expiry     time.Time
	prev, next *Node
}

// DistributedLRUCache represents the LRU cache
type DistributedLRUCache struct {
	mu         sync.RWMutex
	capacity   int
	data       map[string]*Node
	head, tail *Node
	ttl        time.Duration
	consul     *api.Client
	cacheNodes []string
	quit       chan struct{}
	wg         sync.WaitGroup
}

// NewDistributedLRUCache creates a new LRU cache with a given capacity and TTL for items
func NewDistributedLRUCache(capacity int, ttl time.Duration, consulAddr string) (*DistributedLRUCache, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr
	consulClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	kvPairs, _, err := consulClient.KV().List("cache_nodes", nil)
	if err != nil {
		return nil, err
	}

	cacheNodes := make([]string, len(kvPairs))
	for i, kv := range kvPairs {
		cacheNodes[i] = string(kv.Value)
	}

	cache := &DistributedLRUCache{
		capacity:   capacity,
		data:       make(map[string]*Node),
		ttl:        ttl,
		consul:     consulClient,
		cacheNodes: cacheNodes,
		quit:       make(chan struct{}),
	}
	cache.head = &Node{}
	cache.tail = &Node{}
	cache.head.next = cache.tail
	cache.tail.prev = cache.head
	cache.wg.Add(1)
	go cache.expire()
	return cache, nil
}

// getCacheNode retrieves the appropriate cache node using consistent hashing
func (c *DistributedLRUCache) getCacheNode(key string) (string, error) {
	if len(c.cacheNodes) == 0 {
		return "", fmt.Errorf("no cache nodes available")
	}
	index := int(hashKey(key)) % len(c.cacheNodes)
	return c.cacheNodes[index], nil
}

// hashKey provides a simple hash function for consistency
func hashKey(key string) uint64 {
	hash := 0
	for _, b := range key {
		hash = 31*hash + int(b)
	}
	return uint64(hash)
}

// Set adds an item to the cache
func (c *DistributedLRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, ok := c.data[key]
	if ok {
		// Move the existing node to the head
		c.moveToHead(node)
	} else {
		// Add a new node
		newNode := &Node{
			key:    key,
			value:  value,
			expiry: time.Now().Add(c.ttl),
			next:   c.head.next,
			prev:   c.head,
		}
		c.data[key] = newNode
		c.head.next.prev = newNode
		c.head.next = newNode
	}

	// Evict the least recently used item if the cache exceeds capacity
	if len(c.data) > c.capacity {
		oldestNode := c.tail.prev
		delete(c.data, oldestNode.key)
		c.remove(oldestNode)
	}

	// Store the value in the appropriate distributed cache node
	cacheNode, err := c.getCacheNode(key)
	if err != nil {
		fmt.Println("Error retrieving cache node for key", key)
		return
	}
	// Here, we would store the value in the remote cache node `cacheNode`.
	fmt.Printf("Caching key %s in node %s\n", key, cacheNode)
}

// Get retrieves an item from the cache
func (c *DistributedLRUCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	node, ok := c.data[key]
	if !ok {
		return nil, false
	}

	// Check if the item has expired
	if time.Now().After(node.expiry) {
		// Remove expired item
		delete(c.data, key)
		c.remove(node)
		return nil, false
	}

	// Move the accessed item to the head to mark it as recently used
	c.moveToHead(node)
	return node.value, true
}

// remove removes a node from the linked list
func (c *DistributedLRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead moves a node to the head of the linked list (most recently used)
func (c *DistributedLRUCache) moveToHead(node *Node) {
	c.remove(node)
	node.next = c.head.next
	node.prev = c.head
	c.head.next.prev = node
	c.head.next = node
}

// expire periodically checks and removes expired items from the cache
func (c *DistributedLRUCache) expire() {
	defer c.wg.Done()
	for {
		select {
		case <-time.After(c.ttl / 2):
			// Expiration logic is already handled in Get and Set.
		case <-c.quit:
			return
		}
	}
}

// Close shuts down the cache and waits for any cleanup tasks to finish
func (c *DistributedLRUCache) Close() {
	close(c.quit)
	c.wg.Wait()
}

func main() {
	cache, err := NewDistributedLRUCache(2, time.Second*2, "127.0.0.1:8500")
	if err != nil {
		fmt.Println("Error initializing distributed cache:", err)
		return
	}
	defer cache.Close()

	// Set items in the cache
	cache.Set("mykey1", "Hello World!")
	cache.Set("mykey2", "Another Value")

	// Get items from the cache
	value, found := cache.Get("mykey1")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}

	// Access another item to make it the most recently used
	value, found = cache.Get("mykey2")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}

	// Set a third item to trigger eviction
	cache.Set("mykey3", "New Value")

	// The first item should be evicted and not found
	value, found = cache.Get("mykey1")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}
}
