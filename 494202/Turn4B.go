package main

import (
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/logging"
	"context"
)

var (
	logger *logging.Logger
)

func init() {
	// Initialize the logger
	ctx := context.Background()
	logClient, err := logging.NewClient(ctx, "your-project-id")
	if err != nil {
		fmt.Printf("Error creating logging client: %v\n", err)
		return
	}
	logger = logClient.Logger("cloud-function-logger")
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

// LoggingMiddleware logs each request
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		end := time.Now()
		latency := end.Sub(start)

		// Log the request details using GCP Logging
		reqLog := &logging.Entry{
			Payload: map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"user_agent": r.UserAgent(),
				"latency":    latency.String(),
			},
			Severity: logging.Info,
		}
		logger.Log(reqLog)
	}
}

// AddTimeout adds a timeout to the request context
func AddTimeout(duration time.Duration) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()
			next(w, r.WithContext(ctx))
		}
	}
}

// GoogleCloudFunction is the main entry point for the Cloud Function
func GoogleCloudFunction(w http.ResponseWriter, r *http.Request) {
	// The actual logic of the cloud function
	time.Sleep(2 * time.Second) // Simulate work with a delay
	fmt.Fprintf(w, "Hello, Cloud Function! Processed by instance: %s", r.Header.Get("X-Appengine-Instance"))
}

func main() {
	// Define the middleware chain
	handler := Chain(GoogleCloudFunction, AddTimeout(5*time.Second), LoggingMiddleware)
	http.HandleFunc("/", handler)
	fmt.Println("Starting server...")
	http.ListenAndServe(":8080", nil)
}
