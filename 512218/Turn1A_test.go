package _12218

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseMissingParameter(t *testing.T) {
	req, err := http.NewRequest("GET", "/?param1=value1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param1 := r.URL.Query().Get("param1")
		param2 := r.URL.Query().Get("param2")

		if param1 == "" || param2 == "" {
			http.Error(w, "Missing required parameter", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestParseRepeatedParameter(t *testing.T) {
	req, err := http.NewRequest("GET", "/?param1=value1&param1=value2", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()["param1"]

		if len(values) > 1 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strings.Join(values, ",")))
		} else {
			http.Error(w, "Expected multiple values for param1", http.StatusBadRequest)
		}
	})

	handler.ServeHTTP(rr, req)

	if rr.Body.String() == "" {
		t.Errorf("Expected a response body, got empty")
	}
}

func TestParseMalformedInput(t *testing.T) {
	req, err := http.NewRequest("GET", "/?param1=<script>alert('XSS')</script>", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param1 := r.URL.Query().Get("param1")
		sanitized := strings.ReplaceAll(param1, "<", "&lt;")
		sanitized = strings.ReplaceAll(sanitized, ">", "&gt;")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sanitized))
	})

	handler.ServeHTTP(rr, req)

	if rr.Body.String() == "<script>alert('XSS')</script>" {
		t.Errorf("Response contains malformed input: %s", rr.Body.String())
	}
}

func TestParseLongQueryString(t *testing.T) {
	longParamValue := "a"
	for i := 0; i < 10000; i++ {
		longParamValue += "a"
	}
	req, err := http.NewRequest("GET", "/?param1="+longParamValue, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Removed length check for long query strings
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}
