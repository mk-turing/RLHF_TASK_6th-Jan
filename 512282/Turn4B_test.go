package _12282

import (
	"math/rand"
	"runtime"
	"testing"
)

func generateRandomData(size int, maxValue int) []int {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(maxValue)
	}
	return data
}

func simpleInline(x int) int {
	return x + 1
}

func complexInline(x int) int {
	sum := 0
	for i := 0; i < x; i++ {
		for j := 0; j < i; j++ {
			sum += i + j
		}
	}
	return sum
}

func nonInlinedFunction(x int) int {
	return x * x
}

func slidingWindowAnalysis(data []int, windowSize int, funcF func(int) int) []int {
	result := make([]int, len(data)-windowSize+1)
	for i := 0; i < len(result); i++ {
		sum := 0
		for j := i; j < i+windowSize; j++ {
			sum += funcF(data[j])
		}
		result[i] = sum
	}
	return result
}

func BenchmarkSlidingWindowAnalysis(b *testing.B) {
	const dataSize = 100000
	const windowSize = 100
	data := generateRandomData(dataSize, 100)

	benchmarkFuncs := []struct {
		name       string
		funcF      func(int) int
		memorySize int
	}{
		{
			name:       "SimpleInline",
			funcF:      simpleInline,
			memorySize: int(runtime.MemStats{}.Alloc),
		},
		{
			name:       "ComplexInline",
			funcF:      complexInline,
			memorySize: int(runtime.MemStats{}.Alloc),
		},
		{
			name:       "NonInlined",
			funcF:      nonInlinedFunction,
			memorySize: int(runtime.MemStats{}.Alloc),
		},
	}

	b.Run("Baseline", func(b *testing.B) {
		for _, bf := range benchmarkFuncs {
			b.Run(bf.name, func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					slidingWindowAnalysis(data, windowSize, bf.funcF)
				}
			})
		}
	})
}

func TestMemoryUsage(t *testing.T) {
	t.Log("Benchmarking memory usage with different inlined functions...")

	data := generateRandomData(100000, 100)
	windowSize := 100

	funcF := func(f func(int) int) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		slidingWindowAnalysis(data, windowSize, f)
		runtime.ReadMemStats(&ms)
		t.Logf("Memory usage (after analysis): %d bytes\n", ms.Alloc)
	}

	t.Run("SimpleInline", func(t *testing.T) {
		t.Log("Simple inlined function")
		funcF(simpleInline)
	})

	t.Run("ComplexInline", func(t *testing.T) {
		t.Log("Complex inlined function")
		funcF(complexInline)
	})

	t.Run("NonInlined", func(t *testing.T) {
		t.Log("Non-inlined function")
		funcF(nonInlinedFunction)
	})
}
