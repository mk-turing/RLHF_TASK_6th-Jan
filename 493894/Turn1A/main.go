package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"./mock"
	"./service"
)

func Test_BenchmarkServiceProcess(t *testing.T) {
	// Configure mocking
	mockStorage := NewMockStorage(t)
	mockData := []byte("Hello World")
	mockStorage.On("GetData").Return(mockData, nil)

	// Create service with mocked storage
	svc := service.NewService(mockStorage)

	t.Benchmark("MockedServiceProcess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := svc.Process()
			if err != nil {
				b.Fatalf("Service processing failed: %v", err)
			}
		}
	})
}

func Test_BenchmarkRealServiceProcess(t *testing.T) {
	// Mock setup using actual storage (assuming you have some `RealStorage`)
	// MockStorage is used here for simplicity and should be replaced with actual data fetches.

	// Create service with mocked storage (for demonstration)
	svc := service.NewService(mock.NewMockStorage(t))

	// The real implementation involves using non-mocked storage
	// Make changes to this function to use an actual implementation.
	t.Benchmark("RealServiceProcess", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := svc.Process()
			if err != nil {
				b.Fatalf("Service processing failed: %v", err)
			}
		}
	})
}

func main() {
	defer runtime.GC()
	fmt.Printf("%s\n", runBenchmarks())
}

func runBenchmarks() string {
	rand.Seed(time.Now().UnixNano())

	var output bytes.Buffer
	b := &testing.B{}
	b.SetWriter(&output)
	b.Run("", func() {
		t := testing.T{}
		t.Parallel()
		testing.RunParallel(func(pb *testing.PB) {
			t.Run("MockedServiceProcess", func(t *testing.T) {
				s := b.StartTimer("MockedServiceProcess")
				for pb.Next() {
					runBench(t, s)
				}
				s.Stop()
			})
		}, testing.DefaultParallelism)
	})
	return output.String()
}

func runBench(t *testing.T, s *testing.T) {
	// Ensure the mocked data is reset for the next run
	mockStorage := NewMockStorage(t)
	mockData := []byte("Hello World")
	mockStorage.On("GetData").Return(mockData, nil)

	svc := service.NewService(mockStorage)

	_, err := svc.Process()
	if err != nil {
		t.Fatalf("Service processing failed: %v", err)
	}
}
