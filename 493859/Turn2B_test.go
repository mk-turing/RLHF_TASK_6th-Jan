package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

// Dependency interface for swappable dependencies
type Dependency interface {
	Execute() error
}

// BenchmarkResult holds the result of a single benchmark configuration
type BenchmarkResult struct {
	Config string
	Duration time.Duration
	Error error
}

// Registry holds registered dependencies and their configurations
type Registry struct {
	mu           sync.Mutex
	dependencies map[string]Dependency
	activeKey    string
}

// NewRegistry creates a new dependency registry
func NewRegistry() *Registry {
	return &Registry{
		dependencies: make(map[string]Dependency),
	}
}

// Register adds a new dependency implementation with a specific key
func (r *Registry) Register(key string, dep Dependency) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dependencies[key]; exists {
		log.Fatalf("Dependency key %q is already registered", key)
	}
	r.dependencies[key] = dep
}

// SetActive switches the active dependency based on the given key
func (r *Registry) SetActive(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dependencies[key]; !exists {
		log.Fatalf("Dependency key %q is not registered", key)
	}
	r.activeKey = key
}

// Active retrieves the current active dependency
func (r *Registry) Active() Dependency {
	r.mu.Lock()
	defer r.mu.Unlock()
	if dep, exists := r.dependencies[r.activeKey]; exists {
		return dep
	}
	log.Fatalf("No active dependency set")
	return nil
}

// RunBenchmark runs a single benchmark for the active dependency
func (r *Registry) RunBenchmark(b *testing.B, wg *sync.WaitGroup) {
	defer wg.Done()
	dep := r.Active()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dep.Execute(); err != nil {
			b.Fatalf("Error: %v", err)
		}
	}
}

// SlowDependency simulates a slow dependency
type SlowDependency struct{}

func (d *SlowDependency) Execute() error {
	time.Sleep(200 * time.Millisecond)
	return nil
}

// FastDependency simulates a fast dependency
type FastDependency struct{}

func (d *FastDependency) Execute() error {
	return nil
}

// TestDependencies concurrently benchmarks multiple configurations
func TestDependencies(t *testing.T) {
	registry := NewRegistry()
	registry.Register("slow", &SlowDependency{})
	registry.Register("fast", &FastDependency{})

	numConfigs := 3
	results := make(chan BenchmarkResult, numConfigs)
	wg := sync.WaitGroup{}

	// Concurrently run benchmarks for each configuration
	for i := 0; i < numConfigs; i++ {
		wg.Add(1)
		go func(id int) {
			config := fmt.Sprintf("config-%d", id)
			registry.SetActive(config)
			b := testing.Benchmark(func(b *testing.B) {
				registry.RunBenchmark(b, &wg)
			})
			results <- BenchmarkResult{
				Config:   config,
				Duration: b.T,
			}
		}(i)
	}

	// Wait for all benchmarks to complete
	wg.Wait()
	close(results)

	// Aggregate and report the results
	aggregatedResults := make(map[string]time.Duration)