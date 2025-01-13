package main

import (
	"fmt"
	"sync"
)

// Define the Event struct
type Event struct {
	Key   string
	Value string
}

// Define the ShardedMap struct
type ShardedMap struct {
	shards    []map[string]string
	numShards int
	mu        sync.RWMutex
}

// NewShardedMap function creates a new sharded map
func NewShardedMap(numShards int) *ShardedMap {
	shards := make([]map[string]string, numShards)
	for i := range shards {
		shards[i] = make(map[string]string)
	}
	return &ShardedMap{
		shards:    shards,
		numShards: numShards,
	}
}

// Get the hash value for the key
func (sm *ShardedMap) getShardIndex(key string) int {
	// Use a simple hash function for demonstration purposes.
	// In a real application, you would use a more robust hashing algorithm.
	return int(key[0] % uint8(sm.numShards))
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

// Simulate event traffic
func simulateEventTraffic(shardedMap *ShardedMap) {
	for i := 0; i < 10000; i++ {
		event := Event{
			Key:   fmt.Sprintf("key_%d", i),
			Value: fmt.Sprintf("value_%d", i),
		}
		shardedMap.Put(event)
	}
	// Read events
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key_%d", i)
		value, ok := shardedMap.Get(key)
		if !ok {
			fmt.Printf("Error: Key %s not found\n", key)
		}
		if value != fmt.Sprintf("value_%d", i) {
			fmt.Printf("Error: Value for key %s is incorrect\n", key)
		}
	}
}
func main() {
	numShards := 4 // Adjust the number of shards as needed
	shardedMap := NewShardedMap(numShards)
	simulateEventTraffic(shardedMap)
	fmt.Println("Event traffic simulated successfully")
}
