package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

const (
	NumIterations    = 100000
	ProfileInterval = 100 * time.Millisecond
	BufferSize      = 1000
)

type BufferPool struct {
	sync.Pool
}

func (bp *BufferPool) Get() *bytes.Buffer {
	v := bp.Pool.Get()
	if v == nil {
		return bytes.NewBuffer(make([]byte, 0, BufferSize))
	}
	return v.(*bytes.Buffer)
}

func (bp *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.Pool.Put(b)
}

func TestMemoryAllocationOptimized(t *testing.T) {
	bp := &BufferPool{}
	var result int
	for i := 0; i < NumIterations; i++ {
		buffer := bp.Get()
		for j := 0; j < BufferSize; j++ {
			buffer.WriteByte(byte(j))
		}
		result += buffer.Len()
		bp.Put(buffer)
	}
	t.Logf("Sum of buffer lengths: %d", result)
}

func runBenchmark(testFunc func(t *testing.T)) {
	// Start CPU profiling
	cpuprofile := fmt.Sprintf("cpuprofile.%d", time.Now().UnixNano())
	f, err := os.Create(cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Start heap profiling
	heapprofile := fmt.Sprintf("heapprofile.%d", time.Now().UnixNano())
	heapDump := make(chan bool)
	go func() {
		for range time.Tick(ProfileInterval) {
			// Write heap profile to a file
			f, err := os.Create(heapprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.Lookup("heap").WriteTo(f, 0)
			f.Close()
		}
	}()

	// Run the benchmark multiple times to get stable results
	for i := 0; i < NumIterations; i++ {
		testFunc(&testing.T{})
	}

	// Stop heap profiling
	close(heapDump)
}

func generateReports() {
	// Generate CPU profile report
	cmd := exec.Command("go", "tool", "pprof", "cpuprofile.report", "cpuprofile.*")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("CPU profile report generated")

	// Generate heap profile report
	cmd = exec.Command("go", "tool", "pprof", "heapprofile.*", "heapprofile.report")
	var outHeap bytes.Buffer
	cmd.Stdout = &outHeap
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Heap profile report generated")
	fmt.Println(outHeap.String())
}

func main() {
	// Run the benchmark with memory pool
	runBenchmark(TestMemoryAllocationOptimized)
	fmt.Println("Benchmark with memory pool completed.")

	// Run the benchmark without memory pool for comparison
	// Run the benchmark with memory pool
	runBenchmark(TestMemoryAllocationUnoptimized)
	fmt.Println("Benchmark without memory pool completed.")
	// Generate reports