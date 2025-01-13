package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Define the Event struct
type Event struct {
	Key   string
	Value string
}

// Define the TokenBucket struct to implement rate limiting
type TokenBucket struct {
	capacity       int32
	current        int32
	replenishRate  int
	lastRefillTime int64
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity int, replenishRate int) *TokenBucket {
	return &TokenBucket{
		capacity:       int32(capacity),
		current:        int32(capacity),
		replenishRate:  replenishRate,
		lastRefillTime: time.Now().UnixNano(),
	}
}

// CanConsume checks if a token is available
func (tb *TokenBucket) CanConsume() bool {
	return atomic.LoadInt32(&tb.current) > 0
}

// Consume consumes a token
func (tb *TokenBucket) Consume() {
	if !tb.CanConsume() {
		tb.refill()
	}
	atomic.SubtractInt32(&tb.current, 1)
}

// refill replenishes tokens
func (tb *TokenBucket) refill() {
	now := time.Now().UnixNano()
	elapsed := now - tb.lastRefillTime
	replenishment := int(elapsed) / int(time.Second/time.Nanosecond) * tb.replenishRate

	if replenishment > 0 {
		newCurrent := atomic.AddInt32(&tb.current, int32(replenishment))
		if newCurrent > tb.capacity {
			atomic.StoreInt32(&tb.current, tb.capacity)
		}
		tb.lastRefillTime = now
	}
}

// Define the ReplicatedShardedMap struct with rate limiting
type ReplicatedShardedMap struct {
	chr       *ConsistentHashRing
	shards    map[string]map[string]string
	numShards int
	numNodes  int
	replicas  int
	rateLimit map[string]*TokenBucket
	mu        sync.RWMutex
}

// NewReplicatedShardedMap creates a new replicated sharded map with rate limiting
func NewReplicatedShardedMap(numShards int, replicas int, rateLimitCapacity int, rateLimitRefillRate int) *ReplicatedShardedMap {
	shards := make(map[string]map[string]string, numShards)
	rateLimit := make(map[string]*TokenBucket, numShards)
	for i := 0; i < numShards; i++ {
		shards[fmt.Sprintf("%d", i)] = make(map[string]string)
		rateLimit[fmt.Sprintf("%d", i)] = NewTokenBucket(rateLimitCapacity, rateLimitRefillRate)
	}
	return &ReplicatedShardedMap{
		chr:       NewConsistentHashRing(replicas),
		shards:    shards,
		numShards: numShards,
		replicas:  replicas,
		rateLimit: rateLimit,
	}
}

// hashKey is a helper to use the hashKey from ConsistentHashRing
func (rsm *ReplicatedShardedMap) hashKey(key string) uint32 {
	return rsm.chr.hashKey(key)
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

	rsm.chr.RemoveNode(node)
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
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	tb := rsm.rateLimit[fmt.Sprintf("%d", index)]
	tb.Consume()

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
	rsm.mu.Lock()
	defer rsm.mu.Unlock()

	tb := rsm.rateLimit[fmt.Sprintf("%d", index)]
	tb.Consume()

	delete(rsm.shards[fmt.Sprintf("%d", index)], key)

	// Sync data with replicas
	replicas := rsm.getReplicas(primaryNode, index)
	for _, replicaNode := range replicas {
		rsm.syncData(fmt.Sprintf("%d", index), replicaNode)
	}
}

// Simulate event traffic
func simulateEventTraffic(replicatedMap *ReplicatedShardedMap) {
	for i := 0; i < 10000; i++ {
		event := Event{
			Key:   fmt.Sprintf("key_%d", i),
			Value: fmt.Sprintf("value_%d", i),
		}
		replicatedMap.Put(event)
	}
	fmt.Println("Event traffic simulated successfully")
}

// Simulate node failure
func simulateNodeFailure(replicatedMap *ReplicatedShardedMap, node string) {
	replicatedMap.RemoveNode(node)
	fmt.Printf("Node %s failed\n", node)
}

func main() {
	numShards := 4
	replicas := 2
	rateLimitCapacity := 100
	rateLimitRefillRate := 10 // tokens per second
	replicatedMap := NewReplicatedShardedMap(numShards, replicas, rateLimitCapacity, rateLimitRefillRate)

	// Add nodes to the consistent hash ring
	replicatedMap.chr.AddNode("node1")
	replicatedMap.chr.AddNode("node2")
	replicatedMap.chr.AddNode("node3")
	replicatedMap.chr.AddNode("node4")

	simulateEventTraffic(replicatedMap)

	// Simulate node failure to demonstrate failure handling
	time.Sleep(1 * time.Second) // Wait for operations to settle
	simulateNodeFailure(replicatedMap, "node1")

	time.Sleep(1 * time.Second) // Give some time for nodes to adjust and sync

	fmt.Println("Event traffic and node failure simulated successfully")
}
