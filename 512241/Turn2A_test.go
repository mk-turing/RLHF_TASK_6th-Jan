package _12241

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// A/B test function interface
type ABTest interface {
	Run() float64
}

// A/B test implementation
type MyABTest struct {
	count int
}

func (t *MyABTest) Run() float64 {
	// Simulate different workloads or operations
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return float64(time.Now().UnixNano()) / 1e6 // Return execution time in milliseconds
}

// Performance data buffer
type performanceBuffer struct {
	data     []float64
	size     int
	capacity int
	mu       sync.RWMutex
	index    int
}

func newPerformanceBuffer(capacity int) *performanceBuffer {
	return &performanceBuffer{
		data:     make([]float64, capacity),
		size:     0,
		capacity: capacity,
	}
}

func (pb *performanceBuffer) Add(value float64) {
	pb.mu.Lock()
	defer pb.mu.Unlock()
	if pb.size == pb.capacity {
		pb.index = (pb.index + 1) % pb.capacity
	} else {
		pb.size++
	}
	pb.data[pb.index] = value
}

func (pb *performanceBuffer) Average() float64 {
	pb.mu.RLock()
	defer pb.mu.RUnlock()
	if pb.size == 0 {
		return 0
	}
	return sum(pb.data[:pb.size]) / float64(pb.size)
}

func sum(fs []float64) float64 {
	var total float64
	for _, f := range fs {
		total += f
	}
	return total
}

// A/B testing framework with dynamic reallocation
func RunDynamicABTest(tests []ABTest, iterations int, analysisInterval time.Duration) (result map[string]float64) {
	result = make(map[string]float64)
	numTests := len(tests)
	weights := make([]float64, numTests)
	for i := range weights {
		weights[i] = 1.0 / float64(numTests)
	}

	buf := make([]*performanceBuffer, numTests)
	for i := range buf {
		buf[i] = newPerformanceBuffer(100) // Buffer size of 100
	}

	wg := &sync.WaitGroup{}

	// Start analysis goroutine
	go func() {
		for {
			time.Sleep(analysisInterval)
			var bestIndex int
			var minAvg float64 = float64(^uint(0))
			for i, _ := range tests {
				avg := buf[i].Average()
				if avg < minAvg {
					bestIndex = i
					minAvg = avg
				}
			}

			// Reallocate traffic
			for i, w := range weights {
				if i != bestIndex {
					w *= 0.9 // Decrease other weights
				}
			}
			weights[bestIndex] *= 1.1 // Increase best weight
		}
	}()

	// Start test goroutines
	for i := range tests {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				execTime := tests[index].Run()
				buf[index].Add(execTime)
				if rand.Float64() < weights[index] {
					fmt.Printf("Executed variation %d: %.2f ms\n", index, execTime)
				}
			}
		}(i)
	}

	wg.Wait()

	// Calculate average execution time for each test
	for i, _ := range tests {
		avg := buf[i].Average()
		result[fmt.Sprintf("variation%d", i)] = avg
	}

	return result
}

// Benchmark function
func BenchmarkDynamicABTest(b *testing.B) {
	iterations := 1000
	analysisInterval := time.Second / 10

	tests := []ABTest{&MyABTest{count: 1}, &MyABTest{count: 2}}

	for i := 0; i < b.N; i++ {
		result := RunDynamicABTest(tests, iterations, analysisInterval)
		fmt.Printf("Iterations: %d, Result: %v\n", iterations, result)
	}
}
