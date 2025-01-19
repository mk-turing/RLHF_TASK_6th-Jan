package main

import (
	"fmt"
	"net/http"
)

// Middleware type defines the signature for middleware functions
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain allows chaining multiple middleware functions together
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// Iterate in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

// LoggingMiddleware logs each request
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request received: %s %s\n", r.Method, r.URL.Path)
		next(w, r)
	}
}

// AuthMiddleware checks for a mock authentication token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken != "Bearer securetoken" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		next(w, r)
	}
}

// GoogleCloudFunction is the main entry point for the Cloud Function
func GoogleCloudFunction(w http.ResponseWriter, r *http.Request) {
	// The actual logic of the cloud function
	w.Write([]byte("Hello, Cloud Function!"))
}

// EntryPoint wraps the GoogleCloudFunction with all necessary middleware
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	handler := Chain(GoogleCloudFunction, LoggingMiddleware, AuthMiddleware)
	handler(w, r)
}

func main() {
	http.HandleFunc("/", EntryPoint)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
