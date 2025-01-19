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
	Execute() (time.Duration, error)
}

// Registry holds registered dependencies and their configurations.
type Registry struct {
	mu           sync.Mutex
	dependencies map[string]Dependency
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

// Get retrieves a dependency by key.
func (r *Registry) Get(key string) Dependency {
	r.mu.Lock()
	defer r.mu.Unlock()
	if dep, exists := r.dependencies[key]; exists {
		return dep
	}
	log.Fatalf("Dependency key %q is not registered", key)
	return nil
}

// SlowDependency simulates a slow dependency.
type SlowDependency struct{}

func (d *SlowDependency) Execute() (time.Duration, error) {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("SlowDependency executed")
	return time.Since(start), nil
}

// FastDependency simulates a fast dependency.
type FastDependency struct{}

func (d *FastDependency) Execute() (time.Duration, error) {
	start := time.Now()
	fmt.Println("FastDependency executed")
	return time.Since(start), nil
}

// MediumDependency simulates a medium-speed dependency.
type MediumDependency struct{}

func (d *MediumDependency) Execute() (time.Duration, error) {
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	fmt.Println("MediumDependency executed")
	return time.Since(start), nil
}

func TestConcurrentBenchmarking(t *testing.T) {
	registry := NewRegistry()
	registry.Register("slow", &SlowDependency{})
	registry.Register("medium", &MediumDependency{})
	registry.Register("fast", &FastDependency{})

	benchmarks := []string{"slow", "medium", "fast"}
	results := make(map[string]time.Duration)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, benchType := range benchmarks {
		wg.Add(1)
		go func(bench string) {
			defer wg.Done()
			dep := registry.Get(bench)

			var totalDuration time.Duration
			for i := 0; i < 10; i++ { // Run each benchmark 10 times for better averaging
				duration, err := dep.Execute()
				if err != nil {
					log.Fatalf("Error executing %s: %v", bench, err)
				}
				totalDuration += duration
			}
			averageDuration := totalDuration / 10

			mu.Lock()
			results[bench] = averageDuration
			mu.Unlock()
		}(benchType)
	}

	wg.Wait()

	fmt.Println("Benchmark Results:")
	for benchType, duration := range results {
		fmt.Printf("%s: %v\n", benchType, duration)
	}
}
