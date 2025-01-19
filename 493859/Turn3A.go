package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Dependency is the generic interface that dependencies must implement.
type Dependency interface {
	Name() string
	Execute() (time.Duration, error)
}

// DependencyRegistry manages dependency configurations and their implementations.
type DependencyRegistry struct {
	mu           sync.RWMutex
	dependencies map[string]Dependency
}

// NewRegistry creates and returns a new DependencyRegistry instance.
func NewRegistry() *DependencyRegistry {
	return &DependencyRegistry{
		dependencies: make(map[string]Dependency),
	}
}

// Register adds a new dependency to the registry.
func (r *DependencyRegistry) Register(key string, dep Dependency) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.dependencies[key]; exists {
		log.Fatalf("Dependency key %q is already registered", key)
	}
	r.dependencies[key] = dep
}

// LoadDependency dynamically loads a dependency from the registry by its key.
func (r *DependencyRegistry) LoadDependency(key string) (Dependency, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if dep, exists := r.dependencies[key]; exists {
		return dep, nil
	}
	return nil, fmt.Errorf("dependency key %q is not registered", key)
}

// UnloadDependency removes a dependency from the registry.
func (r *DependencyRegistry) UnloadDependency(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.dependencies, key)
}

// SlowDependency simulates a slow dependency.
type SlowDependency struct{}

func (d *SlowDependency) Name() string { return "SlowDependency" }
func (d *SlowDependency) Execute() (time.Duration, error) {
	start := time.Now()
	time.Sleep(100 * time.Millisecond) // Simulate work
	return time.Since(start), nil
}

// FastDependency simulates a fast dependency.
type FastDependency struct{}

func (d *FastDependency) Name() string { return "FastDependency" }
func (d *FastDependency) Execute() (time.Duration, error) {
	start := time.Now()
	return time.Since(start), nil
}

// MediumDependency simulates a moderate-speed dependency.
type MediumDependency struct{}

func (d *MediumDependency) Name() string { return "MediumDependency" }
func (d *MediumDependency) Execute() (time.Duration, error) {
	start := time.Now()
	time.Sleep(50 * time.Millisecond) // Simulate work
	return time.Since(start), nil
}

// CustomDependency allows for custom dependency overriding.
type CustomDependency struct{}

func (d *CustomDependency) Name() string { return "CustomDependency" }
func (d *CustomDependency) Execute() (time.Duration, error) {
	start := time.Now()
	time.Sleep(60 * time.Millisecond) // Custom behavior
	return time.Since(start), nil
}

// runBenchmark runs the benchmark concurrently on the provided configurations.
func runBenchmark(depsRegistry *DependencyRegistry, configs []string) map[string]time.Duration {
	results := make(map[string]time.Duration)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, configKey := range configs {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			if dep, err := depsRegistry.LoadDependency(key); err == nil {
				var totalDuration time.Duration
				for i := 0; i < 10; i++ {
					duration, err := dep.Execute()
					if err != nil {
						fmt.Printf("Error executing %s: %v\n", key, err)
						return
					}
					totalDuration += duration
				}
				averageDuration := totalDuration / 10

				mu.Lock()
				results[dep.Name()] = averageDuration
				mu.Unlock()
			}
		}(configKey)
	}

	wg.Wait()
	return results
}

func main() {
	// Create the dependency registry
	registry := NewRegistry()

	// Registering default dependencies
	registry.Register("slow", &SlowDependency{})
	registry.Register("medium", &MediumDependency{})
	registry.Register("fast", &FastDependency{})

	// Configurations to dynamically benchmark
	configs := []string{"slow", "medium", "fast"}

	// Running initial benchmarks
	results := runBenchmark(registry, configs)
	fmt.Println("Benchmark results:")
	for depName, avgDuration := range results {
		fmt.Printf("%s: %v ms\n", depName, avgDuration.Milliseconds())
	}

	// Unload and load a custom implementation
	registry.UnloadDependency("medium")

	// CustomDependency overrides the default implementation for a specific configuration
	registry.Register("medium", &CustomDependency{})

	// Rerun benchmarks to include the custom implementation
	updatedResults := runBenchmark(registry, configs)
	fmt.Println("Updated benchmark results with custom implementation:")
	for depName, avgDuration := range updatedResults {
		fmt.Printf("%s: %v ms\n", depName, avgDuration.Milliseconds())
	}
}
