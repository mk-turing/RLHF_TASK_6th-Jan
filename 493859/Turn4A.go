package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Dependency interface defines methods for dependencies.
type Dependency interface {
	Name() string
	Execute() (time.Duration, error)
	Initialize() error
	Cleanup()
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

// LogError logs an error message and simulates alerting.
func LogError(message string) {
	log.Printf("Error: %s", message)
	// Simulate alerting stakeholders
	fmt.Printf("ALERT: %s\n", message)
}

// ExecuteWithFaultTolerance executes a dependency with fault tolerance.
func ExecuteWithFaultTolerance(dep Dependency) (time.Duration, error) {
	if err := dep.Initialize(); err != nil {
		LogError(fmt.Sprintf("Initialization failed for %s: %v", dep.Name(), err))
		return 0, err
	}
	defer dep.Cleanup()

	duration, err := dep.Execute()
	if err != nil {
		LogError(fmt.Sprintf("Execution failed for %s: %v", dep.Name(), err))
		return 0, err
	}

	return duration, nil
}

// SlowDependency simulates a slow dependency.
type SlowDependency struct{}

func (d *SlowDependency) Name() string { return "SlowDependency" }
func (d *SlowDependency) Initialize() error {
	// Simulate potential initialization failure
	if time.Now().Unix()%2 == 0 {
		return fmt.Errorf("random initialization failure")
	}
	fmt.Println("SlowDependency initialized")
	return nil
}
func (d *SlowDependency) Execute() (time.Duration, error) {
	start := time.Now()
	time.Sleep(100 * time.Millisecond) // Simulate work
	return time.Since(start), nil
}
func (d *SlowDependency) Cleanup() {
	fmt.Println("SlowDependency cleaned up")
}

// FastDependency simulates a fast dependency.
type FastDependency struct{}

func (d *FastDependency) Name() string { return "FastDependency" }
func (d *FastDependency) Initialize() error {
	fmt.Println("FastDependency initialized")
	return nil
}
func (d *FastDependency) Execute() (time.Duration, error) {
	start := time.Now()
	return time.Since(start), nil
}
func (d *FastDependency) Cleanup() {
	fmt.Println("FastDependency cleaned up")
}

// MediumDependency simulates a moderate-speed dependency.
type MediumDependency struct{}

func (d *MediumDependency) Name() string { return "MediumDependency" }
func (d *MediumDependency) Initialize() error {
	fmt.Println("MediumDependency initialized")
	return nil
}
func (d *MediumDependency) Execute() (time.Duration, error) {
	start := time.Now()
	time.Sleep(50 * time.Millisecond) // Simulate work
	return time.Since(start), nil
}
func (d *MediumDependency) Cleanup() {
	fmt.Println("MediumDependency cleaned up")
}

// Rollback simulates the rollback process.
func Rollback() {
	fmt.Println("System rollback to a stable state")
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
				duration, err := ExecuteWithFaultTolerance(dep)
				if err != nil {
					Rollback()
					return
				}

				mu.Lock()
				results[dep.Name()] = duration
				mu.Unlock()
			}
		}(configKey)
	}

	wg.Wait()
	return results
}

func main() {
	// Set up logging to a file
	logFile, err := os.OpenFile("benchmark.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Create the dependency registry
	registry := NewRegistry()

	// Register default dependencies
	registry.Register("slow", &SlowDependency{})
	registry.Register("medium", &MediumDependency{})
	registry.Register("fast", &FastDependency{})

	// Configurations to dynamically benchmark
	configs := []string{"slow", "medium", "fast"}

	// Running benchmarks with fault tolerance
	results := runBenchmark(registry, configs)
	fmt.Println("Benchmark results:")
	for depName, avgDuration := range results {
		fmt.Printf("%s: %v ms\n", depName, avgDuration.Milliseconds())
	}
}
