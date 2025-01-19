package main

import (
	"flag"
	"fmt"
	"testing"
	"time"
)

type BenchmarkDependency interface {
	RunBenchmark()
}
type MyDependencyA struct{}

func (a *MyDependencyA) RunBenchmark() {
	time.Sleep(10 * time.Millisecond)
}

type MyDependencyB struct{}

func (b *MyDependencyB) RunBenchmark() {
	time.Sleep(20 * time.Millisecond)
}

var (
	dependencyFlag      = flag.String("dependency", "a", "Select the benchmark dependency (a|b)")
	iterationsFlag      = flag.Int("iterations", 1000, "Number of iterations for benchmark")
	benchmarkDependency BenchmarkDependency
)

func init() {
	flag.Parse()
	switch *dependencyFlag {
	case "a":
		benchmarkDependency = &MyDependencyA{}
	case "b":
		benchmarkDependency = &MyDependencyB{}
	default:
		flag.Usage()
	}
}
func BenchmarkFunction(b *testing.B) {
	for i := 0; i < *iterationsFlag; i++ {
		benchmarkDependency.RunBenchmark()
	}
}

func main() {
	fmt.Println("Running benchmark with dependency:", *dependencyFlag)
	testing.Benchmark(BenchmarkFunction)
}
