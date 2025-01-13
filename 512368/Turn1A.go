package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	id   string
	data map[string]interface{}
}

type ShardManager struct {
	nodes []*Node
	wg    *sync.WaitGroup
}

func NewShardManager(nodes int) *ShardManager {
	var sm ShardManager
	sm.nodes = make([]*Node, nodes)
	sm.wg = &sync.WaitGroup{}
	for i := range sm.nodes {
		sm.nodes[i] = &Node{id: fmt.Sprintf("node-%d", i+1), data: make(map[string]interface{})}
	}
	return &sm
}

func (sm *ShardManager) SimulateEvents(numEvents int) {
	sm.wg.Add(numEvents)
	rand.Seed(time.Now().UnixNano())

	go func() {
		for i := 0; i < numEvents; i++ {
			key := fmt.Sprintf("key-%d", rand.Intn(100000))
			value := fmt.Sprintf("value-%d", i)
			sm.UpdateKey(key, value)
			time.Sleep(time.Millisecond * 100)
			sm.wg.Done()
		}
	}()
}

func (sm *ShardManager) UpdateKey(key string, value interface{}) {
	hashValue := sm.hash(key)
	index := hashValue % len(sm.nodes)
	sm.nodes[index].data[key] = value
	fmt.Printf("Key %s inserted into node %s\n", key, sm.nodes[index].id)
}

func (sm *ShardManager) Wait() {
	sm.wg.Wait()
}

func (sm *ShardManager) hash(key string) int {
	hashValue := 0
	for _, byte := range key {
		hashValue = 31*hashValue + int(byte)
	}
	return hashValue
}

func main() {
	sm := NewShardManager(5)
	defer sm.Wait()

	sm.SimulateEvents(100)
	fmt.Println("Event processing completed.")

	for _, node := range sm.nodes {
		fmt.Printf("Node %s Data: %v\n", node.id, node.data)
	}
}
