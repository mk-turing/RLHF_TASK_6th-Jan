package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

// Define the RPC interface
type KVStore interface {
	Get(key string) (string, error)
	Put(key string, value string) error
	Delete(key string) error
}

// Data structure to store key-value pairs
type Store struct {
	data map[string]string
	lock sync.Mutex
}

// RPC handler for Get method
func (s *Store) Get(key string, reply *string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.data[key] == "" {
		return fmt.Errorf("Key not found: %s", key)
	}
	return nil
}

// RPC handler for Put method
func (s *Store) Put(args *struct{ Key, Value string }, reply *string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[args.Key] = args.Value
	*reply = "OK"
	return nil
}

// RPC handler for Delete method
func (s *Store) Delete(key string, reply *string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
	*reply = "OK"
	return nil
}

// Main function to start the RPC server
func main() {
	// Create a new Store instance
	store := &Store{data: make(map[string]string)}

	// Register the RPC server
	rpc.Register(store)
	rpc.HandleHTTP()

	// Listen on a specified TCP port
	l, e := net.Listen("tcp", ":8080")
	if e != nil {
		fmt.Println("Listen error:", e)
		return
	}

	// Serve RPC requests
	fmt.Println("Server listening on :8080")
	rpc.Accept(l)
}
