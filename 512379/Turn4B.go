package main

import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"
	"time"
)

// Node represents a node in the distributed cache cluster
type Node struct {
	id   string
	mu   sync.RWMutex
	data map[string]*NodeEntry
}

// NodeEntry represents an entry in a node's cache
type NodeEntry struct {
	value  interface{}
	expiry time.Time
}

// NewNode creates a new node with a given ID
func NewNode(id string) *Node {
	return &Node{
		id:   id,
		data: make(map[string]*NodeEntry),
	}
}

// ConsistentHashCache represents a distributed consistent hashing cache
type ConsistentHashCache struct {
	mu          sync.RWMutex
	nodes       []*Node
	nodeIDs     []uint32
	replicas    int
	capacity    int
	ttl         time.Duration
	quit        chan struct{}
	wg          sync.WaitGroup
	expireTicker *time.Ticker
}

// NewConsistentHashCache creates a new distributed consistent hashing cache with a given number of nodes and replicas
func NewConsistentHashCache(nodes []*Node, replicas int, capacity int, ttl time.Duration) *ConsistentHashCache {
	c := &ConsistentHashCache{
		nodes:      nodes,
		replicas:   replicas,
		capacity:   capacity,
		ttl:        ttl,
		quit:       make(chan struct{}),
		expireTicker: time.NewTicker(ttl / 2),
	}
	c.wg.Add(1)
	go c.expire()
	return c
}

// hash calculates the hash value of a key
func (c *ConsistentHashCache) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

// getNode finds the node responsible for a given key
func (c *ConsistentHashCache) getNode(key string) *Node {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hash := c.hash(key)
	idx := sort.Search(len(c.nodeIDs), func(i int) bool {
		return c.nodeIDs[i] >= hash
	})
	if idx == len(c.nodeIDs) {
		idx = 0
	}
	return c.nodes[idx]
}

// Set adds an item to the cache
func (c *ConsistentHashCache) Set(key string, value interface{}) {
	node := c.getNode(key)
	node.mu.Lock()
	defer node.mu.Unlock()

	// Update the entry in the local node's cache
	node.data[key] = &NodeEntry{
		value:  value,
		expiry: time.Now().Add(c.ttl),
	}

	// Replicate the entry to replicas
	for i := 1; i <= c.replicas; i++ {
		replicaNode := c.getNode(fmt.Sprintf("%s-replica-%d", key, i))
		replicaNode.mu.Lock()
		replicaNode.data[key] = &NodeEntry{
			value:  value,
			expiry: time.Now().Add(c.ttl),
		}
		replicaNode.mu.Unlock()
	}
}

// Get retrieves an item from the cache
func (c *ConsistentHashCache) Get(key string) (interface{}, bool) {
	node := c.getNode(key)
	node.mu.RLock()