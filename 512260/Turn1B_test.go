package _12260

import (
	"testing"
)

func square(x int) int {
	return x * x
}

func BenchmarkSquare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		square(1000)
	}
}

func TestSquare(t *testing.T) {
	result := square(5)
	if result != 25 {
		t.Errorf("Expected 25, got %d", result)
	}
}

func main() {
}
