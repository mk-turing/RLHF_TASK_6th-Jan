package main

import (
	"fmt"
	"testing"
)

// BenchmarkMyFunction benchmarks the performance of MyFunction with increasing input size.
func BenchmarkMyFunction(b *testing.B) {
	for _, size := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				//MyFunction(createRandomData(size))
			}
		})
	}
}
