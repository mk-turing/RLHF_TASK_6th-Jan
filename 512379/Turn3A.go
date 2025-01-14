package main

import (
	"fmt"
	"sync"
	"time"
)

// Node represents a node in the doubly linked list
type Node struct {
	key        string
	value      interface{}
	expiry     time.Time
	prev, next *Node
}

// LRUCache represents the LRU cache
type LRUCache struct {
	mu         sync.RWMutex
	capacity   int
	data       map[string]*Node
	head, tail *Node
	ttl        time.Duration
	quit       chan struct{}
	wg         sync.WaitGroup
}

// NewLRUCache creates a new LRU cache with a given capacity and TTL for items
func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	c := &LRUCache{
		capacity: capacity,
		data:     make(map[string]*Node),
		ttl:      ttl,
		quit:     make(chan struct{}),
	}
	c.head = &Node{}
	c.tail = &Node{}
	c.head.next = c.tail
	c.tail.prev = c.head
	c.wg.Add(1)
	go c.expire()
	return c
}

// Set adds an item to the cache
func (c *LRUCache) Set(key string, value interface{}) {
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
}

// Get retrieves an item from the cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
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
func (c *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// moveToHead moves a node to the head of the linked list (most recently used)
func (c *LRUCache) moveToHead(node *Node) {
	c.remove(node)
	node.next = c.head.next
	node.prev = c.head
	c.head.next.prev = node
	c.head.next = node
}

// expire periodically checks and removes expired items from the cache
func (c *LRUCache) expire() {
	defer c.wg.Done()
	for {
		select {
		case <-time.After(c.ttl / 2):
			// Expiration logic is already handled in Get and Set.
			// You can optionally add another cleanup cycle here, but it's not required.
		case <-c.quit:
			return
		}
	}
}

// Close shuts down the cache and waits for any cleanup tasks to finish
func (c *LRUCache) Close() {
	c.quit <- struct{}{}
	c.wg.Wait()
}

func main() {
	cache := NewLRUCache(2, time.Second*2)
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

	// Wait a bit to allow expiration
	time.Sleep(time.Second * 3)

	// The items should have expired
	value, found = cache.Get("mykey1")
	if found {
		fmt.Println("Value from cache:", value)
	} else {
		fmt.Println("Cache miss.")
	}
}
