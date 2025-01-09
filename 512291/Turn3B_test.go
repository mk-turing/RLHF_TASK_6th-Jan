package _12291

import (
	"runtime"
	"testing"
)

type Shape interface {
	Area() float64
}

type Square struct {
	side float64
}

func (s Square) Area() float64 {
	return s.side * s.side
}

type Rectangle struct {
	length, width float64
}

func (r Rectangle) Area() float64 {
	return r.length * r.width
}

// Calculation functions for different parameter types and concurrency models

func calculateAreaDirect(side float64) float64 {
	return side * side
}

func calculateAreaConcurrent(side float64, results chan<- float64) {
	results <- side * side
}

var sq = Square{side: 100.0}
var rc = Rectangle{length: 100.0, width: 200.0}
var sf = 100.0

func BenchmarkSquareDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sq.Area()
	}
}

func BenchmarkRectangleDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rc.Area()
	}
}

func BenchmarkSquareConcurrent(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	results := make(chan float64, b.N)
	for i := 0; i < b.N; i++ {
		go calculateAreaConcurrent(sf, results)
	}
	for i := 0; i < b.N; i++ {
		<-results
	}
}
