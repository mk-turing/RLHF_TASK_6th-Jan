package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// Define the Event struct
type Event struct {
	Key   string
	Value string
}

// Define the ConsistentHashRing struct
type ConsistentHashRing struct {
	ring     map[uint32][]string
	replicas int
	nodes    []string
	mu       sync.RWMutex
}

// NewConsistentHashRing creates a new hash ring
func NewConsistentHashRing(replicas int) *ConsistentHashRing {
	return &ConsistentHashRing{
		ring:     make(map[uint32][]string),
		replicas: replicas,
		nodes:    []string{},
	}
}

// Get the hash value for a key
func (chr *ConsistentHashRing) hashKey(key string) uint32 {
	h := fnv.New32()
	h.Write([]byte(key))
	return h.Sum32()
}

// Add a node to the ring
func (chr *ConsistentHashRing) AddNode(node string) {
	chr.mu.Lock()
	defer chr.mu.Unlock()

	chr.nodes = append(chr.nodes, node)
	for i := 0; i < chr.replicas; i++ {
		vnode := fmt.Sprintf("%s-%d", node, i)
		hash := chr.hashKey(vnode)
		chr.ring[hash] = append(chr.ring[hash], vnode)
	}
}

// Remove a node from the ring
func (chr *ConsistentHashRing) RemoveNode(node string) {
	chr.mu.Lock()
	defer chr.mu.Unlock()

	// Remove virtual nodes for the specified node
	for i := 0; i < chr.replicas; i++ {
		vnode := fmt.Sprintf("%s-%d", node, i)
		hash := chr.hashKey(vnode)
		if _, ok := chr.ring[hash]; ok {
			// Remove the vnode from the slice
			newSlice := []string{}
			for _, v := range chr.ring[hash] {
				if v != vnode {
					newSlice = append(newSlice, v)
				}
			}
			if len(newSlice) == 0 {
				delete(chr.ring, hash)
			} else {
				chr.ring[hash] = newSlice
			}
		}
	}

	// Remove the node from the node list
	var newNodes []string
	for _, n := range chr.nodes {
		if n != node {
			newNodes = append(newNodes, n)
		}
	}
	chr.nodes = newNodes
}

// GetNode retrieves the node responsible for a given key
func (chr *ConsistentHashRing) GetNode(key string) string {
	chr.mu.RLock()
	defer chr.mu.RUnlock()

	hash := chr.hashKey(key)
	for i, vnode := range chr.nodes {
		vnodeHash := chr.hashKey(vnode)
		if vnodeHash > hash {
			return vnode
		}
		if i == len(chr.nodes)-1 {
			return chr.nodes[0] // Wrap around to the first node
		}
	}
	return chr.nodes[0] // In case of no nodes
}

// Define the ShardedMap struct
type ShardedMap struct {
	chr       *ConsistentHashRing
	shards    map[string]map[string]string
	numShards int
	mu        sync.RWMutex
}

// NewShardedMap creates a new sharded map
func NewShardedMap(numShards int, replicas int) *ShardedMap {
	shards := make(map[string]map[string]string, numShards)
	for i := 0; i < numShards; i++ {
		shards[fmt.Sprintf("%d", i)] = make(map[string]string)
	}
	return &ShardedMap{
		chr:       NewConsistentHashRing(replicas),
		shards:    shards,
		numShards: numShards,
	}
}

// hashKey is a helper to use the hashKey from ConsistentHashRing
func (sm *ShardedMap) hashKey(key string) uint32 {
	return sm.chr.hashKey(key)
}

// Put adds or updates a key-value pair
func (sm *ShardedMap) Put(event Event) {
	node := sm.chr.GetNode(event.Key)
	index := sm.hashKey(node) % uint32(sm.numShards)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.shards[fmt.Sprintf("%d", index)][event.Key] = event.Value
}

// Get retrieves the value for a given key
func (sm *ShardedMap) Get(key string) (string, bool) {
	node := sm.chr.GetNode(key)
	index := sm.hashKey(node) % uint32(sm.numShards)
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.shards[fmt.Sprintf("%d", index)][key]
	return value, ok
}

// Remove deletes a key-value pair
func (sm *ShardedMap) Remove(key string) {
	node := sm.chr.GetNode(key)
	index := sm.hashKey(node) % uint32(sm.numShards)
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.shards[fmt.Sprintf("%d", index)], key)
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
	numShards := 4
	replicas := 3
	shardedMap := NewShardedMap(numShards, replicas)

	// Add nodes to the consistent hash ring
	shardedMap.chr.AddNode("node1")
	shardedMap.chr.AddNode("node2")
	shardedMap.chr.AddNode("node3")

	simulateEventTraffic(shardedMap)
	fmt.Println("Event traffic simulated successfully")
}
