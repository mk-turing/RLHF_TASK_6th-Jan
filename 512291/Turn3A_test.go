package _12291

import (
	"sync"
	"testing"
)

type Point struct {
	x, y float64
}

func calculateAreaStruct(p Point) float64 {
	return p.x * p.y
}

func calculateAreaArray(a [2]float64) float64 {
	return a[0] * a[1]
}

func calculateAreaSlice(s []float64) float64 {
	return s[0] * s[1]
}

var p = Point{x: 100.0, y: 200.0}
var a = [2]float64{100.0, 200.0}
var s = []float64{100.0, 200.0}

func BenchmarkCalculateAreaStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaStruct(p)
	}
}

func BenchmarkCalculateAreaArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaArray(a)
	}
}

func BenchmarkCalculateAreaSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaSlice(s)
	}
}

// Benchmark functions with parallel execution

func parallelCalculateArea(b *testing.B, f func(float64, float64) float64, x, y float64, numWorkers int) {
	wg := &sync.WaitGroup{}
	results := make(chan float64, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < b.N/numWorkers; j++ {
				results <- f(x, y)
			}
		}()
	}

	wg.Wait()
	close(results)

	// Consume all results
	for range results {
	}
}

func BenchmarkParallelCalculateAreaStruct(b *testing.B) {
	parallelCalculateArea(b, calculateAreaStruct, p.x, p.y, b.N/1000)
}

func BenchmarkParallelCalculateAreaArray(b *testing.B) {
	parallelCalculateArea(b, calculateAreaArray, a[0], a[1], b.N/1000)
}

func BenchmarkParallelCalculateAreaSlice(b *testing.B) {
	parallelCalculateArea(b, calculateAreaSlice, s[0], s[1], b.N/1000)
}
