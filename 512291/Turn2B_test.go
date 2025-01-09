package _12291

import (
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

// Adding benchmarks for function call strategies with different parameter types:

func calculateAreaStructPointer(p *Point) float64 {
	return p.x * p.y
}

func calculateAreaArrayPointer(a *[2]float64) float64 {
	return a[0] * a[1]
}

func calculateAreaSlicePointer(s *[]float64) float64 {
	return (*s)[0] * (*s)[1]
}

var pPtr = &p
var aPtr = &a
var sPtr = &s

func BenchmarkCalculateAreaStructPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaStructPointer(pPtr)
	}
}

func BenchmarkCalculateAreaArrayPointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaArrayPointer(aPtr)
	}
}

func BenchmarkCalculateAreaSlicePointer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaSlicePointer(sPtr)
	}
}
