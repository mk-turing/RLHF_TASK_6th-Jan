package _12218

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSQLInjectionQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "/?username=admin'--", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		// Prevent SQL Injection by sanitizing the input
		sanitizedUsername := strings.ReplaceAll(username, "'", "''")
		sanitizedUsername = strings.ReplaceAll(sanitizedUsername, "--", "") // Remove SQL comments

		sqlQuery := "SELECT * FROM users WHERE username = '" + sanitizedUsername + "'"

		// Simulate executing the SQL query
		fmt.Fprintf(w, "Query: %s\n", sqlQuery)
	})

	handler.ServeHTTP(rr, req)

	expectedQuery := "Query: SELECT * FROM users WHERE username = 'admin'''\n"
	if rr.Body.String() != expectedQuery {
		t.Errorf("Expected sanitized query: %s, got: %s", expectedQuery, rr.Body.String())
	}
}

func TestScriptInjectionQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "/?comment=<script>alert('XSS')</script>", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		comment := r.URL.Query().Get("comment")
		// Prevent script injection by sanitizing the input
		sanitizedComment := strings.ReplaceAll(comment, "<", "&lt;")
		sanitizedComment = strings.ReplaceAll(sanitizedComment, ">", "&gt;")
		sanitizedComment = strings.ReplaceAll(sanitizedComment, "'", "&#x27;") // Encode single quotes

		fmt.Fprintf(w, "Comment: %s\n", sanitizedComment)
	})

	handler.ServeHTTP(rr, req)

	expectedComment := "Comment: &lt;script&gt;alert(&#x27;XSS&#x27;)&lt;/script&gt;\n"
	if rr.Body.String() != expectedComment {
		t.Errorf("Expected sanitized comment: %s, got: %s", expectedComment, rr.Body.String())
	}
}

func TestDoSWithLongQueryString(t *testing.T) {
	longParamValue := "a"
	for i := 0; i < 10000; i++ {
		longParamValue += "a"
	}

	for i := 0; i < 1000; i++ { // Simulate rapid requests
		req, err := http.NewRequest("GET", "/?param1="+longParamValue, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate handling a large request
			w.WriteHeader(http.StatusOK)
			time.Sleep(10 * time.Millisecond) // Deliberate sleep to simulate resource usage
		})

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	}
}

type RateLimiter struct {
	maxRequests int
	windowTime  time.Duration
	mutex       sync.Mutex
	requests    []time.Time
}

func NewRateLimiter(maxRequests int, windowTime time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		windowTime:  windowTime,
		requests:    make([]time.Time, 0),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.windowTime)

	// Remove expired requests
	i := 0
	for _, t := range rl.requests {
		if t.After(windowStart) {
			break
		}
		i++
	}
	rl.requests = rl.requests[i:]

	// Check if request is allowed
	if len(rl.requests) < rl.maxRequests {
		rl.requests = append(rl.requests, now)
		return true
	}

	return false
}

func TestRateLimiting(t *testing.T) {
	maxRequests := 5
	maxTime := time.Second

	rateLimiter := NewRateLimiter(maxRequests, maxTime)

	for i := 0; i < 10; i++ { // Simulate exceeding the rate limit
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rateLimiter.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}

			w.WriteHeader(http.StatusOK)
		})

		handler.ServeHTTP(rr, req)

		if i < maxRequests && rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d before limit, got %d", http.StatusOK, rr.Code)
		}

		if i >= maxRequests && rr.Code != http.StatusTooManyRequests {
			t.Errorf("Expected status code %d after limit, got %d", http.StatusTooManyRequests, rr.Code)
		}
	}
}
