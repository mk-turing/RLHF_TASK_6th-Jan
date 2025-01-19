package main

import (
	"fmt"
	"log"
	"net/http"
)

// Logger interface for abstracting logging implementation
type Logger interface {
	Printf(format string, v ...interface{})
}

// Middleware type defines the signature for middleware functions
type Middleware func(Logger, http.HandlerFunc) http.HandlerFunc

// Chain allows chaining multiple middleware functions together
func Chain(l Logger, f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// Iterate in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](l, f)
	}
	return f
}

// LoggingMiddleware logs each request
func LoggingMiddleware(l Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				l.Printf("Error in LoggingMiddleware: %v", err)
			}
		}()

		l.Printf("Request received: %s %s\n", r.Method, r.URL.Path)
		next(w, r)
	}
}

// AuthMiddleware checks for a mock authentication token
func AuthMiddleware(l Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				l.Printf("Error in AuthMiddleware: %v", err)
			}
		}()

		authToken := r.Header.Get("Authorization")
		if authToken != "Bearer securetoken" {
			l.Printf("Unauthorized request: %s", r.URL.Path)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		next(w, r)
	}
}

// GoogleCloudFunction is the main entry point for the Cloud Function
func GoogleCloudFunction(l Logger, w http.ResponseWriter, r *http.Request) error {
	// The actual logic of the cloud function
	// Simulate an error
	if r.URL.Path == "/error" {
		return fmt.Errorf("internal server error")
	}

	w.Write([]byte("Hello, Cloud Function!"))
	return nil
}

// EntryPoint wraps the GoogleCloudFunction with all necessary middleware
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	l := log.New(w, "", 0) // Using the response writer for logging in this example
	handler := Chain(l, GoogleCloudFunction, LoggingMiddleware, AuthMiddleware)

	if err := handler(w, r); err != nil {
		l.Printf("Error in EntryPoint: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", EntryPoint)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
