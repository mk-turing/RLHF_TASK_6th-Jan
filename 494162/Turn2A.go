package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

// Define the service interface
type KVService interface {
	Set(key, value string, version int64, newVersion *int64) error
	Get(key string, value *string, version *int64) error
	Delete(key string, version int64) error
}

// Implement the service
type KV struct {
	store map[string]versionedKeyValue
	mu    sync.Mutex
}

type versionedKeyValue struct {
	Value   string
	Version int64
}

func (kv *KV) Set(key, value string, version int64, newVersion *int64) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if existing, ok := kv.store[key]; ok && existing.Version > version {
		// Conflict: the current version on server is higher
		*newVersion = existing.Version
		return fmt.Errorf("conflict: current version is %d", existing.Version)
	}

	// Update the key with the new value and increment the version
	kv.store[key] = versionedKeyValue{value, version + 1}
	*newVersion = version + 1
	return nil
}

func (kv *KV) Get(key string, value *string, version *int64) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if vkv, ok := kv.store[key]; ok {
		*value = vkv.Value
		*version = vkv.Version
		return nil
	}
	*value = ""
	*version = 0
	return fmt.Errorf("key not found")
}

func (kv *KV) Delete(key string, version int64) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	existing, ok := kv.store[key]
	if ok && existing.Version == version {
		delete(kv.store, key)
		return nil
	}
	return fmt.Errorf("conflict: current version is %d", existing.Version)
}

func main() {
	// Create an instance of the service
	kv := &KV{store: make(map[string]versionedKeyValue)}

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
