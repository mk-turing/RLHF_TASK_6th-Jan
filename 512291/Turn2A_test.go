package _12291

import (
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

func calculateAreaDirect(side float64) float64 {
	return side * side
}

func calculateAreaDirectArray(sides []float64) float64 {
	sum := 0.0
	for _, side := range sides {
		sum += side * side
	}
	return sum
}

func calculateAreaDirectSlice(sides []float64) float64 {
	sum := 0.0
	for _, side := range sides {
		sum += side * side
	}
	return sum
}

func calculateAreaWithClosure(side float64) float64 {
	return calculateAreaClosure(side)()
}

func calculateAreaClosure(side float64) func() float64 {
	return func() float64 {
		return side * side
	}
}

var sq = Square{side: 100.0}
var rc = Rectangle{length: 100.0, width: 200.0}
var sf = 100.0
var cf = calculateAreaClosure(100.0)
var sidesArray [100]float64
var sidesSlice []float64

func init() {
	for i := 0; i < 100; i++ {
		sidesArray[i] = float64(i + 1)
		sidesSlice = append(sidesSlice, float64(i+1))
	}
}

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

func BenchmarkSquareInterface(b *testing.B) {
	var s Shape = sq
	for i := 0; i < b.N; i++ {
		s.Area()
	}
}

func BenchmarkRectangleInterface(b *testing.B) {
	var s Shape = rc
	for i := 0; i < b.N; i++ {
		s.Area()
	}
}

func BenchmarkCalculateAreaDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaDirect(sf)
	}
}

func BenchmarkCalculateAreaDirectArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaDirectArray(sidesArray)
	}
}

func BenchmarkCalculateAreaDirectSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaDirectSlice(sidesSlice)
	}
}

func BenchmarkCalculateAreaClosure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cf()
	}
}
