package _12241

import (
	"fmt"
	"math"
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
	data    []float64
	size    int
	capacity int
	mu      sync.RWMutex
	index   int
}

func newPerformanceBuffer(capacity int) *performanceBuffer {
	return &performanceBuffer{
		data:    make([]float64, capacity),
		size:    0,
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

// Chi-square test for independence
func chiSquareTest(obs []int) float64 {
	rows := len(obs)
	cols := len(obs) / rows

	expected := make([]int, rows*cols)
	total := sumInts(obs)
	for i := range expected {
		expected[i] = total * (rows + cols - 1) / (rows * cols)
	}

	var chiSquare float64
	for i := range obs {
		row := i / cols
		col := i % cols
		observed := obs[i]
		chiSquare += math.Pow(float64(observed)-float64(expected[i]), 2) / float64(expected[i])
	}

	return chiSquare
}

func sumInts(ints []int) int {
	var total int
	for _, i := range ints {
		total += i
	}
	return total
}

// A/B testing framework with statistical analysis
func RunStatisticalABTest(tests []ABTest, iterations int, analysisInterval time.Duration) (result map[string]float64, significance float64) {
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