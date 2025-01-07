package main

import (
	"fmt"
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
