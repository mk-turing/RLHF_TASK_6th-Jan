package _12260

import (
	"testing"
)

func add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	t.Run("basic addition", func(t *testing.T) {
		if add(1, 2) != 3 {
			t.Errorf("Expected 3, got %d", add(1, 2))
		}
	})
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		add(1, 2)
	}
}
