// Sample Mock Server implementation
package _12218

import (
	"fmt"
	"net/http"
	"strings"
)

type MockServer struct {
	routes map[string]func(http.ResponseWriter, *http.Request)
}

func NewMockServer() *MockServer {
	return &MockServer{
		routes: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
}

func (m *MockServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.routes[pattern] = handler
}

func (m *MockServer) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, m)
}

func (m *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := m.routes[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	handler(w, r)
}

func TestQueryParameterParsingInDistributedSystem(t *testing.T) {
	// Start the central mock server
	mockServer := NewMockServer()
	mockServer.ListenAndServe(":8080")

	//  (1)  Tests with conflicting formats
	mockServer.HandleFunc("/microservice1/route", func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("format")
		validFormats := []string{"json", "xml"}
		// (..) Perform validation on param and respond appropriately
	})

	mockServer.HandleFunc("/microservice2/route", func(w http.ResponseWriter, r *http.Request) {
		format := r.URL.Query().Get("format")
		if format != "json" {
			http.Error(w, "Invalid format", http.StatusBadRequest)
		}
		// Continue processing
	})


	//  (2) Tests for cascading failures
	// Create a mock handler that triggers a failure in Microservice 2
	mockServer.HandleFunc("/microservice1/route", func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("invalid_param")
		if param != "" {
			// Simulate failure by responding with an error
			http.Error(w, "Invalid input", http.StatusBadRequest)
		}
		http.Error(w, "Failed to handle request", http.StatusInternalServerError)
	})

	//  (3) Tests for Asynchronous communication
	mockServer.HandleFunc("/queue/microservice2", func(w http.ResponseWriter, r *http.Request) {
		// Read message from the queue and parse query parameters
		msg := r.FormValue("msg")
		// Parse query parameters from msg and perform processing
	})

	// Send requests with various scenarios to the mock server
	// Test for conflicting formats
	res, _ := http.Get("http://localhost:8080/microservice1/route?format=xls")

	// Test for cascading failures
	res, _ = http.Get("http://localhost:8080/microservice1/route?invalid_param=true")