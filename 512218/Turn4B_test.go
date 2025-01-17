package _12218

import (
	"net/http"
	"testing"
)

func TestStatefulMockServer(t *testing.T) {
	ms := NewStatefulMockServer()
	defer ms.Close()

	client := &http.Client{}

	// Test 1: Simulate stateful behavior
	req, err := http.NewRequest("GET", ms.URL()+"?param1=value1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	if res.StatusCode != http.StatusFound {
		t.Errorf("Expected status code %d, got %d", http.StatusFound, res.StatusCode)
	}

	// Test 2: Simulate dynamic edge-case parameter generation
	req, err = http.NewRequest("GET", ms.URL(), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	res, err = client.Do(req)
	if err != nil {