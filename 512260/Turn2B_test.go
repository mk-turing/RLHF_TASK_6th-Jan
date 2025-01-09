package _12260

import (
	"sync"
	"testing"
)

func add(a, b int) int {
	return a + b
}

func BenchmarkAddConcurrent(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go func() {
			add(1, 2)
			wg.Done()
		}()
	}

	wg.Wait()
}
