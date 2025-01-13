package main

import (
	"fmt"
	"sync"
	"time"
)

// Define the Event struct
type Event struct {
	Key   string
	Value string
}

// Define the TokenBucket struct
type TokenBucket struct {
	capacity int
	rate     int
	tokens   int
	lastFill time.Time
	mu       sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity int, rate int) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		rate:     rate,
		tokens:   capacity,
		lastFill: time.Now(),
	}
}

// Take takes a token from the bucket
func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastFill)
	tb.lastFill = now

	newTokens := int(elapsed.Seconds() * float64(tb.rate))
	tb.tokens = min(tb.capacity, tb.tokens+newTokens)

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Define the ReplicatedShardedMap struct
type ReplicatedShardedMap struct {
	chr          *ConsistentHashRing
	shards       map[string]map[string]string
	numShards    int
	numNodes     int
	replicas     int
	tokenBuckets []*TokenBucket
	mu           sync.RWMutex
	counter      int
}

// hashKey is a helper to use the hashKey from ConsistentHashRing
func (rsm *ReplicatedShardedMap) hashKey(key string) uint32 {
	return rsm.hashKey(key)
}

// NewReplicatedShardedMap creates a new replicated sharded map
func NewReplicatedShardedMap(numShards int, replicas int, rate int) *ReplicatedShardedMap {
	shards := make(map[string]map[string]string, numShards)
	tokenBuckets := make([]*TokenBucket, numShards)
	for i := 0; i < numShards; i++ {
		shards[fmt.Sprintf("%d", i)] = make(map[string]string)
		tokenBuckets[i] = NewTokenBucket(1000, rate) // Adjust the capacity as needed
	}
	return &ReplicatedShardedMap{
		chr:          NewConsistentHashRing(replicas),
		shards:       shards,
		numShards:    numShards,
		replicas:     replicas,
		tokenBuckets: tokenBuckets,
	}
}

// Add a node to the ring and update replica logic
func (rsm *ReplicatedShardedMap) AddNode(node string) {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	rsm.chr.AddNode(node)
	rsm.numNodes++

	// Update replica assignment
	for _, nodeId := range rsm.chr.nodes {
		index := rsm.hashKey(nodeId) % uint32(rsm.numShards)
		replicas := rsm.getReplicas(nodeId, index)
		rsm.reassignReplicas(index, replicas)
	}
}

// Remove a node from the ring and update replica logic
func (rsm *ReplicatedShardedMap) RemoveNode(node string) {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	rsm.RemoveNode(node)
	rsm.numNodes--

	// Update replica assignment
	for _, nodeId := range rsm.chr.nodes {
		index := rsm.hashKey(nodeId) % uint32(rsm.numShards)
		replicas := rsm.getReplicas(nodeId, index)
		rsm.reassignReplicas(index, replicas)
	}
}

// Get replica nodes for a specific shard index
func (rsm *ReplicatedShardedMap) getReplicas(primaryNode string, index uint32) []string {
	replicas := make([]string, rsm.replicas)
	count := 0
	for _, nodeId := range rsm.chr.nodes {
		if nodeId == primaryNode {
			continue
		}
		replicaIndex := rsm.hashKey(nodeId) % uint32(rsm.numShards)
		if replicaIndex == index {
			replicas[count] = nodeId
			count++
			if count == rsm.replicas {
				break
			}
		}
	}
	return replicas
}

// Reassign replicas for a shard
func (rsm *ReplicatedShardedMap) reassignReplicas(index uint32, replicas []string) {
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	shardKey := fmt.Sprintf("%d", index)
	for _, nodeId := range replicas {
		rsm.syncData(shardKey, nodeId)
	}
}

// Put adds or updates a key-value pair and syncs with replicas
func (rsm *ReplicatedShardedMap) Put(event Event) {
	primaryNode := rsm.chr.GetNode(event.Key)
	index := rsm.hashKey(primaryNode) % uint32(rsm.numShards)

	// Check if a token is available for the shard
	if !rsm.tokenBuckets[index].Take() {
		return // Drop the event if there is no token
	}

	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	rsm.shards[fmt.Sprintf("%d", index)][event.Key] = event.Value

	// Sync data with replicas
	replicas := rsm.getReplicas(primaryNode, index)
	for _, replicaNode := range replicas {
		rsm.syncData(fmt.Sprintf("%d", index), replicaNode)
	}
}

// Get retrieves the value for a given key
func (rsm *ReplicatedShardedMap) Get(key string) (string, bool) {
	primaryNode := rsm.chr.GetNode(key)
	index := rsm.hashKey(primaryNode) % uint32(rsm.numShards)
	rsm.mu.RLock()
	defer rsm.mu.RUnlock()

	value, ok := rsm.shards[fmt.Sprintf("%d", index)][key]
	return value, ok
}

// Remove deletes a key-value pair and syncs with replicas
func (rsm *ReplicatedShardedMap) Remove(key string) {
	primaryNode := rsm.chr.GetNode(key)
	index := rsm.hashKey(primaryNode) % uint32(rsm.numShards)

	// Check if a token is available for the shard
	if !rsm.tokenBuckets[index].Take() {
		return // Drop the event if there is no token
	}

	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	delete(rsm.shards[fmt.Sprintf("%d", index)], key)

	// Sync data with replicas
	replicas := rsm.getReplicas(primaryNode, index)
	for _, replicaNode := range replicas {
		rsm.syncData(fmt.Sprintf("%d", index), replicaNode)
	}
}

// Simulate node failure
func simulateNodeFailure(replicatedMap *ReplicatedShardedMap, node string) {
	replicatedMap.RemoveNode(node)
	fmt.Printf("Node %s failed\n", node)
}

// Simulate event traffic
func simulateEventTraffic(replicatedMap *ReplicatedShardedMap, rate int) {
	// Simulate events at the specified rate
	ticker := time.NewTicker(time.Duration(1000/rate) * time.Millisecond)
	defer ticker.Stop()

	for i := 0; ; i++ {
		<-ticker.C
		event := Event{
			Key:   fmt.Sprintf("key_%d", i),
			Value: fmt.Sprintf("value_%d", i),
		}
		replicatedMap.Put(event)
	}
}

// syncData handles syncing the shard data between nodes.
// Avoid double-locking the mutex by adjusting the locking mechanism.
func (rsm *ReplicatedShardedMap) syncData(shardKey, node string) {
	// Locking the hash ring is fine, as we're not modifying shard data here.
	rsm.chr.mu.Lock() // Lock for modifying the hash ring
	defer rsm.chr.mu.Unlock()

	// Now handle data synchronization without double-locking the shard map mutex.
	// Assume this is a quick operation and simulate it.
	fmt.Printf("Syncing shard %s with node %s\n", shardKey, node)
	time.Sleep(100 * time.Millisecond) // Simulate sync duration
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

func main() {
	numShards := 4
	replicas := 2
	rate := 10000 // Events per second
	replicatedMap := NewReplicatedShardedMap(numShards, replicas, rate)

	// Add nodes to the consistent hash ring
	replicatedMap.chr.AddNode("node1")
	replicatedMap.chr.AddNode("node2")
	replicatedMap.chr.AddNode("node3")
	replicatedMap.chr.AddNode("node4")

	// Simulate event traffic
	go simulateEventTraffic(replicatedMap, rate)

	// Simulate node failure to demonstrate failure handling
	time.Sleep(1 * time.Second) // Wait for operations to settle
	simulateNodeFailure(replicatedMap, "node1")

	time.Sleep(1 * time.Second) // Give some time for nodes to adjust and sync

	fmt.Println("Event traffic and node failure simulated successfully")
}
