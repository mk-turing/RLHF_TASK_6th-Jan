package _12291

import (
	"fmt"
	"os"
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

func calculateAreaDirect(side float64) float64 {
	return side * side
}

func BenchmarkSquareDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculateAreaDirect(100.0)
	}
}

func main() {
	// Run tests and collect profiles
	testing.Main()

	// Run go tool pprof to analyze the profiles
	fmt.Println("Analyzing CPU profile...")
	if err := os.StartProcess("go", []string{"tool", "pprof", "cpu.prof"}, nil); err != nil {
		fmt.Println("Error analyzing CPU profile:", err)
	}

	fmt.Println("Analyzing memory profile...")
	if err := os.StartProcess("go", []string{"tool", "pprof", "mem.prof"}, nil); err != nil {
		fmt.Println("Error analyzing memory profile:", err)
	}
}
