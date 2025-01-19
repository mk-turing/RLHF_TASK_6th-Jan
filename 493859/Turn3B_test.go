package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

// Dependency interface for swappable components
type Dependency interface {
	Execute() (time.Duration, error)
}

// BenchmarkConfig represents a configuration for a benchmark
type BenchmarkConfig struct {
	Name       string
	Dependency string
	Override   Dependency
}

// DependencyRegistry is a registry for all known dependencies
type DependencyRegistry struct {
	mu           sync.Mutex
	dependencies map[string]Dependency
}

// Register a dependency with a key
func (r *DependencyRegistry) Register(key string, dep Dependency) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.dependencies[key]; ok {
		log.Fatalf("Dependency '%s' already registered", key)
	}
	r.dependencies[key] = dep
}

// Get retrieves a dependency by key
func (r *DependencyRegistry) Get(key string) Dependency {
	r.mu.Lock()
	defer r.mu.Unlock()
	if dep, ok := r.dependencies[key]; ok {
		return dep
	}
	log.Fatalf("Dependency '%s' not found", key)
	return nil
}

// DependencyFactory creates dependencies dynamically
type DependencyFactory func() Dependency

// RegisterFactory registers a factory function for a dependency key
func (r *DependencyRegistry) RegisterFactory(key string, factory DependencyFactory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.dependencies[key]; ok {
		log.Fatalf("Dependency '%s' already registered", key)
	}
	r.dependencies[key] = factory()
}

// NewDependencyRegistry creates a new dependency registry
func NewDependencyRegistry() *DependencyRegistry {
	return &DependencyRegistry{
		dependencies: make(map[string]Dependency),
	}
}

// BenchmarkResult represents the result of a single benchmark
type BenchmarkResult struct {
	Name     string
	Duration time.Duration
	Err      error
}

// RunBenchmark executes a single benchmark configuration
func RunBenchmark(cfg BenchmarkConfig, registry *DependencyRegistry, wg *sync.WaitGroup) {
	defer wg.Done()

	var dep Dependency
	// If the config has an override, use it
	if cfg.Override != nil {
		dep = cfg.Override
	} else {
		dep = registry.Get(cfg.Dependency)
	}

	duration, err := dep.Execute()

	result := BenchmarkResult{
		Name:     cfg.Name,
		Duration: duration,
		Err:      err,
	}

	fmt.Println(result)
}

func TestConcurrentBenchmarking(t *testing.T) {
	registry := NewDependencyRegistry()

	// Register default dependency implementation for "sqrt"
	registry.Register("sqrt", newSqrtDependency())

	// Register factory function for "complexsqrt"
	registry.RegisterFactory("complexsqrt", func() Dependency {
		return newComplexSqrtDependency()
	})

	benchmarks := []BenchmarkConfig{
		{Name: "SimpleSqrt", Dependency: "sqrt"},
		{Name: "ComplexSqrt", Dependency: "complexsqrt"},
		// Override the "simple_pi" dependency with a new implementation for this specific benchmark
		{Name: "TestCircleArea", Dependency: "pi", Override: newSimplePiDependency()},
	}

	var wg sync.WaitGroup
	for _, cfg := range benchmarks {
		wg.Add(1)
		go RunBenchmark(cfg, registry, &wg)
	}
	wg.Wait()
}
