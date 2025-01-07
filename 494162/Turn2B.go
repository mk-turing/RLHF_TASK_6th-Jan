package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Define the service interface
type KVService interface {
	Set(key, value string) error
	Get(key string, value *string) error
	Delete(key string) error
}

type LockService interface {
	AcquireLock(key string, timeout time.Duration) error
	ReleaseLock(key string) error
}

// Implement the lock service
type Lock struct {
	locks sync.Map // Map of locks: key -> *sync.Mutex
}

func (l *Lock) AcquireLock(key string, timeout time.Duration) error {
	// Create a new mutex for the key if it doesn't exist
	mu, _ := l.locks.LoadOrStore(key, &sync.Mutex{})

	// Try to acquire the lock within the timeout period
	done := make(chan struct{})
	go func() {
		mu.(*sync.Mutex).Lock()
		close(done)
	}()

	select {
	case <-done:
		return nil // Lock acquired
	case <-time.After(timeout):
		return errors.New("lock acquisition timed out")
	}
}

func (l *Lock) ReleaseLock(key string) error {
	mu, ok := l.locks.Load(key)
	if !ok {
		return errors.New("lock not found")
	}
	mu.(*sync.Mutex).Unlock()
	return nil
}

// Implement the KV service
type KV struct {
	store map[string]string
	lock  *Lock
}

func (kv *KV) Set(key, value string) error {
	// Acquire lock for the key
	if err := kv.lock.AcquireLock(key, 5*time.Second); err != nil {
		return err
	}
	defer kv.lock.ReleaseLock(key)

	kv.store[key] = value
	return nil
}

func (kv *KV) Get(key string, value *string) error {
	// Acquire lock for the key
	if err := kv.lock.AcquireLock(key, 5*time.Second); err != nil {
		return err
	}
	defer kv.lock.ReleaseLock(key)

	if _, ok := kv.store[key]; ok {
		*value = kv.store[key]
		return nil
	}
	*value = ""
	return fmt.Errorf("key not found")
}

func (kv *KV) Delete(key string) error {
	// Acquire lock for the key
	if err := kv.lock.AcquireLock(key, 5*time.Second); err != nil {
		return err
	}
	defer kv.lock.ReleaseLock(key)

	delete(kv.store, key)
	return nil
}

func main() {
	// Create instances of the services
	kv := &KV{store: make(map[string]string), lock: &Lock{}}
	lockService := &Lock{}

	// Register the services with the RPC server
	rpc.Register(kv)
	rpc.Register(lockService)