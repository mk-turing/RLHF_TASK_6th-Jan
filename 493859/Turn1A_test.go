package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

// Dependency represents a generic interface for swappable dependencies.
type Dependency interface {
	Execute() error
}

// Registry holds registered dependencies and their configurations.
type Registry struct {
	mu           sync.Mutex
	dependencies map[string]Dependency
	activeKey    string
}

// NewRegistry creates a new dependency registry.
func NewRegistry() *Registry {
	return &Registry{
		dependencies: make(map[string]Dependency),
	}
}

// Register adds a new dependency implementation with a specific key.
func (r *Registry) Register(key string, dep Dependency) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dependencies[key]; exists {
		log.Fatalf("Dependency key %q is already registered", key)
	}
	r.dependencies[key] = dep
}

// SetActive switches the active dependency based on the given key.
func (r *Registry) SetActive(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dependencies[key]; !exists {
		log.Fatalf("Dependency key %q is not registered", key)
	}
	r.activeKey = key
}

// Active retrieves the current active dependency.
func (r *Registry) Active() Dependency {
	r.mu.Lock()
	defer r.mu.Unlock()
	if dep, exists := r.dependencies[r.activeKey]; exists {
		return dep
	}
	log.Fatalf("No active dependency set")
	return nil
}

type SlowDependency struct{}

func (d *SlowDependency) Execute() error {
	time.Sleep(100 * time.Millisecond)
	fmt.Println("SlowDependency executed")
	return nil
}

// FastDependency simulates a fast dependency.
type FastDependency struct{}

func (d *FastDependency) Execute() error {
	fmt.Println("FastDependency executed")
	return nil
}

func BenchmarkDependencies(b *testing.B) {
	registry := NewRegistry()
	registry.Register("slow", &SlowDependency{})
	registry.Register("fast", &FastDependency{})

	for _, benchType := range []string{"slow", "fast"} {
		b.Run(benchType, func(b *testing.B) {
			registry.SetActive(benchType)
			dep := registry.Active()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := dep.Execute(); err != nil {
					b.Fatalf("Error: %v", err)
				}
			}
		})
	}
}
