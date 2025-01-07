package _93902

import (
	"testing"
)

// A simple addition operation to benchmark
func benchmarkAdd(b *testing.B) {
	var a = 100000
	var b = 200000

	for i := 0; i < b.N; i++ {
		_ = a + b
	}
}

func TestMain(m *testing.M) {
	testing.Init()
	m.Run()
}
