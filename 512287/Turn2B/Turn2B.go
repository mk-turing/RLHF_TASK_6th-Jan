package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime/pprof"
	"testing"
	"time"
)

const (
	// NumIterations defines the number of iterations for the benchmark
	NumIterations = 10
	// ProfileInterval defines the interval at which CPU profiles will be taken
	ProfileInterval = 100 * time.Millisecond
)

func TestMemoryAllocation(t *testing.T) {
	// Example test code
	// Your test code goes here
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
	// Run the benchmark
	runBenchmark(TestMemoryAllocation)

	// Generate reports
	generateReports()
}
