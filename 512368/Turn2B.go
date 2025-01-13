package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Define the Event struct
type Event struct {
	Key   string
	Value string
}

// Define the ShardedMap struct
type ShardedMap struct {
	shards   []map[string]string
	numShards int
	mu      sync.RWMutex
	nodeIDs  []int
}

// NewShardedMap function creates a new sharded map
func NewShardedMap(numShards int) *ShardedMap {
	shards := make([]map[string]string, numShards)
	for i := range shards {
		shards[i] = make(map[string]string)
	}
	return &ShardedMap{
		shards:   shards,
		numShards: numShards,
		nodeIDs:  []int{0}, // Start with one node
	}
}

// getShardIndex uses consistent hashing to determine the shard index for a given key.
func (sm *ShardedMap) getShardIndex(key string) int {
	// Calculate the hash value for the key
	hash := sm.hash(key)

	// Find the shard index based on the hash value
	i := sort.SearchInts(sm.nodeIDs, hash)
	if i == len(sm.nodeIDs) {
		i = 0
	}
	return i % sm.numShards
}

// hash function using a simple linear congruential generator (LCG) for demonstration purposes.
// In a real application, you would use a more robust hashing algorithm like MurmurHash3 or xxHash.
func (sm *ShardedMap) hash(key string) int {
	const (
		a = 1664525
		c = 1013904223
		m = 1 << 32
	)
	hash := uint32(1)
	for _, b := range []byte(key) {
		hash = a*(hash^uint32(b)) + c
		hash %= m
	}
	return int(hash)
}

// Put function adds or updates a key-value pair
func (sm *ShardedMap) Put(event Event) {
	shardIndex := sm.getShardIndex(event.Key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.shards[shardIndex][event.Key] = event.Value
}

// Get function retrieves the value for a given key
func (sm *ShardedMap) Get(key string) (string, bool) {
	shardIndex := sm.getShardIndex(key)
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.shards[shardIndex][key]
	return value, ok
}

// Remove function deletes a key-value pair
func (sm *ShardedMap) Remove(key string) {
	shardIndex := sm.getShardIndex(key)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.shards[shardIndex], key)
}

// AddNode adds a new node to the cluster
func (sm *ShardedMap) AddNode() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	newNodeID := sm.getAvailableNodeID()
	sm.nodeIDs = append(sm.nodeIDs, newNodeID)
	sort.Ints(sm.nodeIDs)
}

// RemoveNode removes a node from the cluster
func (sm *ShardedMap) RemoveNode() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if len(sm.nodeIDs) > 1 {
		nodeToRemoveIndex := rand.Intn(len(sm.nodeIDs) - 1) + 1 // Avoid removing the only node
		sm.nodeIDs = append(sm.nodeIDs[:nodeToRemoveIndex], sm.nodeIDs[nodeToRemoveIndex+1:]...)
	}
}

// getAvailableNodeID finds an available node ID that is not currently being used.
func (sm *ShardedMap) getAvailableNodeID() int {