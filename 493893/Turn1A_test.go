package main

import (
	"flag"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"testing"
)

func main() {
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	benchmarkName := os.Args[1]
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if t, ok := testing.Benchmarks[benchmarkName]; ok {
		b := &testing.B{
			N:               100, // Number of times to run the benchmark
			Setup:           func(b *testing.B) {},
			TearDown:        func(b *testing.B) {},
			MemAlloc:        func() ([]byte, error) { return nil, nil },
			MemFree:         func(b []byte) {},
			TraceHook:       nil,
			BeforeFork:      func() error { return nil },
			AfterFork:       func() error { return nil },
			Hook:            func(*testing.B) {},
			WriteToTemplate: nil,
			ThreadCpuClock:  func(*testing.B) {},
			NeedToStopGc:    false,
			ParallelWorkload: func() []*testing.B {
				return make([]*testing.B, 0)
			},
			Namespace:        "",
			NeedRacing:       false,
			Replicate:        0,
			UseParameter:     false,
			WithFailure:      false,
			Group:            "",
			RemoveIterations: false,
		}
		*b = *t
		b.Run()
	} else {
		panic("Benchmark not found: " + benchmarkName)
	}
}

// Example benchmark to sort an array of random integers
func BenchmarkRandomSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		n := 1000
		slice := make([]int, n)
		rand.Seed(int64(i))
		for j := 0; j < n; j++ {
			slice[j] = rand.Intn(10000)
		}
		sort.Ints(slice)
	}
}
