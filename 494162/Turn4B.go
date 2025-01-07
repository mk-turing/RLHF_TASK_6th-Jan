package main

import (
	"fmt"
	"github.com/hashicorp/raft"
	"log"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"
)

// Define the service interface
type KVService interface {
	Set(key, value string) error
	Get(key string, value *string) error
	Delete(key string) error
	GetVersion() int64
	Join(serverID, addr string) error
}

// KVStore holds the map data and raft state.
type KVStore struct {
	raft          *raft.Raft
	raftNodeID    raft.NodeID
	raftAddr      string
	raftLogStore  raft.LogStore
	raftSnapStore raft.SnapshotStore
	store         map[string]string
	version       int64
	mu            sync.Mutex
}

// NewKVStore creates a new KVStore with raft config
func NewKVStore(nodeID raft.NodeID, addr string) (*KVStore, error) {
	config := raft.DefaultConfig()
	config.LocalID = nodeID

	// In-memory log and snapshot store for simplicity
	logStore := raft.NewInmemStore()
	snapStore := raft.NewInmemStore()

	kvStore := &KVStore{
		raftNodeID:    nodeID,
		raftAddr:      addr,
		raftLogStore:  logStore,
		raftSnapStore: snapStore,
		store:         make(map[string]string),
	}

	// Create the raft instance
	raft, err := raft.NewRaft(config, kvStore, logStore, snapStore, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	kvStore.raft = raft

	return kvStore, nil
}

// Implement the raft.FSM interface to handle Raft log entries
func (kv *KVStore) Apply(log *raft.Log) interface{} {
	var command kvCommand
	if err := log.Decode(&command); err != nil {
		panic(fmt.Sprintf("failed to decode command: %s", err))
	}

	kv.mu.Lock()
	defer kv.mu.Unlock()
	switch command.Op {
	case "set":
		kv.store[command.Key] = command.Value
	case "delete":
		delete(kv.store, command.Key)
	}
	atomic.AddInt64(&kv.version, 1)
	return nil
}

func (kv *KVStore) Snapshot() (raft.FSMSnapshot, error) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	snapshot := &kvSnapshot{
		store:   make(map[string]string),
		version: kv.version,
	}
	for k, v := range kv.store {
		snapshot.store[k] = v
	}
	return snapshot, nil
}

func (kv *KVStore) Restore(rc io.ReadCloser) error {
	defer rc.Close()
	d := json.NewDecoder(rc)
	var snapshot kvSnapshot
	if err := d.Decode(&snapshot); err != nil {
		return err
	}
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store = snapshot.store
	kv.version = snapshot.version
	return nil
}

// Define the service interface
type KVService interface {
	Set(key, value string) error
	Get(key string, value *string) error
	Delete(key string) error
	GetVersion() int64
}

// Implement the service
type KV struct {
	store   map[string]string
	version int64
	mu      sync.Mutex
}

func (kv *KV) Set(key, value string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = value
	atomic.AddInt64(&kv.version, 1)
	return nil
}

func (kv *KV) Get(key string, value *string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if val, ok := kv.store[key]; ok {
		*value = val
		return nil
	}
	*value = ""
	return fmt.Errorf("key not found")
}

func (kv *KV) Delete(key string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.store, key)
	atomic.AddInt64(&kv.version, 1)
	return nil
}

func (kv *KV) GetVersion() int64 {
	return atomic.LoadInt64(&kv.version)
}

func main() {
	// Create an instance of the service
	kv := &KV{store: make(map[string]string)}

	// Register the service with the RPC server
	rpc.Register(kv)

	// Listen on a port
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Accept connections and handle them with the RPC server
	log.Println("RPC server started on :12345")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
