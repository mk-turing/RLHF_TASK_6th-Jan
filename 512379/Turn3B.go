package main

import (
	"fmt"
	"sync"
	"time"
)

type node struct {
	key   string
	value interface{}
	prev  *node
	next  *node
}

type Cache struct {
	mu       sync.RWMutex
	data     map[string]*node
	head     *node
	tail     *node
	capacity int
	ttl      time.Duration
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewCache(capacity int, ttl time.Duration) *Cache {
	c := &Cache{
		data:     make(map[string]*node),
		capacity: capacity,
		ttl:      ttl,
		quit:     make(chan struct{}),
	}
	c.wg.Add(1)
	go c.expire()
	return c
}

func (c *Cache) moveToHead(n *node) {
	if n == c.head {
		return
	}
	if n == c.tail {
		c.tail = n.prev
		c.tail.next = nil
	} else {
		n.prev.next = n.next
		n.next.prev = n.prev
	}
	n.next = c.head
	n.prev = nil
	c.head.prev = n
	c.head = n
}

func (c *Cache) removeLast() {
	if c.tail == nil {
		return
	}
	if c.tail == c.head {
		c.head = nil
	} else {
		c.tail.prev.next = nil
	}
	delete(c.data, c.tail.key)
	c.tail = c.tail.prev
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.data[key]; ok {
		// Update the value if the key already exists
		n.value = value
		c.moveToHead(n)
	} else {
		// Create a new node and add it to the front of the list
		n := &node{key: key, value: value}
		n.next = c.head
		if c.head != nil {
			c.head.prev = n
		}
		c.head = n
		if c.tail == nil {
			c.tail = n
		}
		c.data[key] = n

		if len(c.data) > c.capacity {
			// Remove the last node if the cache is full
			c.removeLast()
		}
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if n, ok := c.data[key]; ok {
		c.moveToHead(n)
		return n.value, true
	}
	return nil, false
}

func (c *Cache) expire() {
	defer c.wg.Done()
	for {
		select {
		case <-time.After(c.ttl / 2):
			now := time.Now()
			c.mu.Lock()
			for key, n := range c.data {
				if time.Now().After(n.value.(time.Time).Add(c.ttl)) {
					// Remove expired entries
					c.moveToHead(n)
					c.removeLast()
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		case <-c.quit:
			return
		}
	}