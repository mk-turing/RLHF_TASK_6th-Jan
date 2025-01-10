package _12283

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Simulate random network latency (0-100ms)
func addLatency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		latency := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(latency)
		next.ServeHTTP(w, r)
	})
}

// Simulate a network failure
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

func BenchmarkTestServer(b *testing.B) {
	// Create a test server with randomized latencies and network failures
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addLatency(simulateNetworkFailure(http.HandlerFunc(testServerHandler))).ServeHTTP(w, r)
	}))
	defer ts.Close()

	client := &http.Client{Timeout: 10 * time.Second}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", ts.URL, nil)
		_, err := client.Do(req)
		if err != nil {
			b.Error(err)
		}
	}
}

// Benchmark user authentication and database access
func BenchmarkAuthAndDataAccess(b *testing.B) {
	// Simulate user auth and database query
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate auth check
		time.Sleep(100 * time.Microsecond)
		// Simulate database query
		doSomeWork()
	}
}
