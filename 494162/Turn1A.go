package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Define the service interface
type KVService interface {
	Set(key, value string) error
	Get(key string, value *string) error
	Delete(key string) error
}

// Implement the service
type KV struct {
	store map[string]string
}

func (kv *KV) Set(key, value string) error {
	kv.store[key] = value
	return nil
}

func (kv *KV) Get(key string, value *string) error {
	if _, ok := kv.store[key]; ok {
		return nil
	}
	*value = ""
	return fmt.Errorf("key not found")
}

func (kv *KV) Delete(key string) error {
	delete(kv.store, key)
	return nil
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
