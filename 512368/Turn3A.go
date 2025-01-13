package main

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

// Define the Event struct
type Event struct {
	Key   string
	Value string
}

// Define the ReplicatedShardedMap struct
type ReplicatedShardedMap struct {
	chr       *ConsistentHashRing
	shards    map[string]map[string]string
	numShards int
	numNodes  int
	replicas  int
	mu        sync.RWMutex
}

// NewReplicatedShardedMap creates a new replicated sharded map
func NewReplicatedShardedMap(numShards int, replicas int) *ReplicatedShardedMap {
	shards := make(map[string]map[string]string, numShards)
	for i := 0; i < numShards; i++ {
		shards[fmt.Sprintf("%d", i)] = make(map[string]string)
	}
	return &ReplicatedShardedMap{
		chr:       NewConsistentHashRing(replicas), // Initialize the ConsistentHashRing
		shards:    shards,
		numShards: numShards,
		replicas:  replicas,
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

func main() {
	numShards := 4
	replicas := 2
	replicatedMap := NewReplicatedShardedMap(numShards, replicas)

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
