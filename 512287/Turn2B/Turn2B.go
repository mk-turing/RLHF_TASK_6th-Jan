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
	// Example test code: simulating memory allocation
	for i := 0; i < 100000; i++ {
		_ = make([]byte, 1024) // allocating memory
	}
}

func runBenchmark(testFunc func(t *testing.T)) (string, string) {
	// Generate unique filenames for CPU and heap profiles
	cpuprofile := fmt.Sprintf("cpuprofile.%d", time.Now().UnixNano())
	heapprofile := fmt.Sprintf("heapprofile.%d", time.Now().UnixNano())

	// Start CPU profiling
	f, err := os.Create(cpuprofile)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// Start heap profiling
	heapDump := make(chan bool)
	go func() {
		for range time.Tick(ProfileInterval) {
			// Write heap profile to a file
			f, err := os.Create(heapprofile)
			if err != nil {
				log.Fatal(err)
			}
			err = pprof.Lookup("heap").WriteTo(f, 0)
			if err != nil {
				log.Fatal(err)
			}
			f.Close()
		}
	}()

	// Run the benchmark multiple times to get stable results
	for i := 0; i < NumIterations; i++ {
		testFunc(&testing.T{})
	}

	// Stop heap profiling by closing the channel
	close(heapDump)

	return cpuprofile, heapprofile
}

func generateReports(cpuprofile, heapprofile string) {
	// Generate CPU profile report
	cmd := exec.Command("go", "tool", "pprof", cpuprofile, "cpuprofile.report")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("CPU profile report generated:", out.String())

	// Generate heap profile report
	cmd = exec.Command("go", "tool", "pprof", heapprofile, "heapprofile.report")
	var outHeap bytes.Buffer
	cmd.Stdout = &outHeap
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Heap profile report generated:", outHeap.String())
}

func main() {
	// Run the benchmark and get the profile filenames
	cpuprofile, heapprofile := runBenchmark(TestMemoryAllocation)

	// Generate reports using the profile filenames
	generateReports(cpuprofile, heapprofile)
}
