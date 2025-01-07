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
	GetReplicas() []string
}

// Implement the service
type KV struct {
	store    map[string]string
	version  int64
	replicas []string
	mu       sync.Mutex
}

func (kv *KV) Set(key, value string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = value
	atomic.AddInt64(&kv.version, 1)
	return kv.replicate("set", key, value)
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
	return kv.replicate("delete", key, "")
}

func (kv *KV) GetVersion() int64 {
	return atomic.LoadInt64(&kv.version)
}

func (kv *KV) GetReplicas() []string {
	return kv.replicas
}

func (kv *KV) replicate(action string, key, value string) error {
	for _, replica := range kv.replicas {
		client, err := rpc.DialHTTP("tcp", replica)
		if err != nil {
			log.Printf("failed to dial %s: %s", replica, err)
			continue
		}
		defer client.Close()
		var args struct {
			Action string
			Key    string
			Value  string
		}
		args.Action = action
		args.Key = key
		args.Value = value
		var reply struct {
			Error string
		}
		if err := client.Call("KV.Replicate", &args, &reply); err != nil {
			log.Printf("failed to replicate to %s: %s", replica, err)
			return fmt.Errorf("replication failed to %s: %s", replica, err)
		}
		if reply.Error != "" {
			log.Printf("replication error to %s: %s", replica, reply.Error)
			return fmt.Errorf("replication error to %s: %s", replica, reply.Error)
		}
	}
	return nil
}

// Replicate method to handle incoming replication requests
func (kv *KV) Replicate(args *struct {
	Action string
	Key    string
	Value  string
}, reply *struct {
	Error string
}) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	switch args.Action {
	case "set":
		kv.store[args.Key] = args.Value
		atomic.AddInt64(&kv.version, 1)
	case "delete":
		delete(kv.store, args.Key)
		atomic.AddInt64(&kv.version, 1)
	default:
		reply.Error = fmt.Sprintf("unknown action: %s", args.Action)
		return nil
	}
	reply.Error = ""
	return nil
}

func main() {
	// Create an instance of the service with replicas
	kv := &KV{store: make(map[string]string), replicas: []string{"localhost:12346", "localhost:12347"}}

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
