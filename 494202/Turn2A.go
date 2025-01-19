package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

// ErrorHandlingMiddleware captures and manages errors raised by subsequent middlewares or the handler
func ErrorHandlingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logError(fmt.Errorf("internal server error: %v", err))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// LoggingMiddleware logs each request
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Request received: %s %s\n", r.Method, r.URL.Path)
		next(w, r)
	}
}

// AuthMiddleware checks for a mock authorization token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		if authToken != "Bearer securetoken" {
			handleClientError(w, "Unauthorized", http.StatusUnauthorized)
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
	handler := Chain(GoogleCloudFunction, ErrorHandlingMiddleware, LoggingMiddleware, AuthMiddleware)
	handler(w, r)
}

// logError logs error messages to stderr or other logging facilities
func logError(err error) {
	log.Printf("Error: %v", err) // Consider enhancing this with request IDs or metadata if required
}

// handleClientError handles client errors, responds to the client, and can log them if necessary
func handleClientError(w http.ResponseWriter, message string, code int) {
	response := map[string]string{"error": message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
	log.Printf("Client error (%d): %s", code, message)
}

func main() {
	http.HandleFunc("/", EntryPoint)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on :%s\n", port)
	http.ListenAndServe(":"+port, nil)
}
