package _12258

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	inputSizes = []int{10000, 100000, 1000000} // Different input sizes to test
	profileDir = "./profiles"                  // Directory to store memory profiles
)

func init() {
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		log.Fatalf("Error creating profile directory: %v", err)
	}
}

func testMemoryLeak(t *testing.T) {
	for _, size := range inputSizes {
		t.Run(fmt.Sprintf("TestSize%d", size), func(t *testing.T) {
			// Start memory profiling
			f, err := os.Create(fmt.Sprintf("%s/test_size_%d.prof", profileDir, size))
			if err != nil {
				t.Fatalf("Error creating memory profile file: %v", err)
			}
			defer f.Close()
			if err := pprof.StartCPUProfile(f); err != nil {
				t.Fatalf("Error starting CPU profile: %v", err)
			}
			defer pprof.StopCPUProfile()

			// Run the test function with the current input size
			testFunction(size)

			// Force garbage collection
			runtime.GC()

			// Sleep for a short duration to ensure any lingering allocations are cleaned up
			time.Sleep(100 * time.Millisecond)

			// Read the memory profile after the test function
			memProfile, err := readMemProfile()
			if err != nil {
				t.Fatalf("Error reading memory profile: %v", err)
			}

			// Analyze the memory profile for leaks
			leaks, err := analyzeMemProfile(memProfile)
			if err != nil {
				t.Fatalf("Error analyzing memory profile: %v", err)
			}

			if len(leaks) > 0 {
				t.Errorf("Detected memory leaks in test size %d:", size)
				for _, leak := range leaks {
					t.Errorf("- %v", leak)
				}
			}
		})
	}
}

// Your test function goes here
func testFunction(size int) {
	// Allocate memory and perform operations here
	// Simulate memory leaks by not deallocating allocated resources
}

// Reads the memory profile after the test function
func readMemProfile() ([]byte, error) {
	cmd := exec.Command("go", "tool", "pprof", "-text", "-inuse_objects", os.Args[0])
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running pprof: %v", err)
	}
	return out.Bytes(), nil
}

// Analyzes the memory profile for leaks
func analyzeMemProfile(profile []byte) ([]string, error) {
	var leaks []string
	lines := strings.Split(string(profile), "\n")
	for _, line := range lines {
		if strings.Contains(line, "leaked") {
			fields := strings.Fields(line)
			if len(fields) > 3 {
				size, err := strconv.Atoi(fields[1])
				if err != nil {
					return nil, fmt.Errorf("error parsing profile size: %v", err)
				}
				leaks = append(leaks, fmt.Sprintf("%d bytes in %s", size, strings.Join(fields[3:], " ")))
			}
		}
	}
	return leaks, nil
}

func TestMemoryLeak(t *testing.T) {
	testMemoryLeak(t)
}
