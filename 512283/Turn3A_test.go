package _12283

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func addLatency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		latency := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(latency)
		next.ServeHTTP(w, r)
	})
}

func simulateNetworkFailure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(10) == 0 {
			http.Error(w, "Network Failure", http.StatusServiceUnavailable)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func testServerHandler(w http.ResponseWriter, r *http.Request) {
	doSomeWork() // Simulate actual work
	fmt.Fprintf(w, "Hello, World!")
}

func doSomeWork() {
	// Simulate various work items
	switch rand.Intn(3) {
	case 0:
		time.Sleep(1 * time.Millisecond)
	case 1:
		time.Sleep(2 * time.Millisecond)
	case 2:
		time.Sleep(5 * time.Millisecond) // Simulate slower work
	}
}

// BenchmarkTestServer with variable number of goroutines
func BenchmarkTestServerParallel(b *testing.B) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addLatency(simulateNetworkFailure(http.HandlerFunc(testServerHandler))).ServeHTTP(w, r)
	}))
	defer ts.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", ts.URL, nil)
			_, err := client.Do(req)
			if err != nil {
				// Use b.Fatalf to report errors in benchmarks
				b.Fatalf("Error in benchmark: %v", err)
			}
		}
	})
}

// Benchmark user authentication and database access with variable goroutines
func BenchmarkAuthAndDataAccessParallel(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Simulate auth check
			time.Sleep(100 * time.Microsecond)
			// Simulate database query
			doSomeWork()
		}
	})
}

// Benchmark for resource contention
func BenchmarkResourceContention(b *testing.B) {
	var sharedResource int

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sharedResource++
			sharedResource--
		}
	})
}

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())
	m.Run()
}
