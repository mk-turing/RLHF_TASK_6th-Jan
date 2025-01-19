package main

import (
	"context"
	"fmt"
	"net/http"
)

// Middleware is a function type that represents a middleware component.
type Middleware func(http.Handler) http.Handler

// handler wraps the actual cloud function handler with all the middlewares.
func handler(ctx context.Context, req *http.Request) (int, string) {
	// Main cloud function logic here
	return http.StatusOK, "Hello, Cloud Functions!"
}

// LoggingMiddleware logs the request and response details.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request:", r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Println("Response sent:", w.Header().Get("Content-Type"))
	})
}

// AuthenticationMiddleware authenticates the request using a custom strategy (e.g., JWT).
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Custom authentication logic here
		// ...

		// If authentication succeeds, call the next handler
		next.ServeHTTP(w, r)
	})
}

// chainMiddlewares chains together multiple middleware components using the callback.
func chainMiddlewares(next http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		next = m(next)
	}
	return next
}

// main function exports the cloud function entry point.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Configure your middleware components here
		middlewares := []Middleware{LoggingMiddleware, AuthenticationMiddleware}

		// Chain the middlewares and call the actual handler
		chainMiddlewares(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			status, message := handler(r.Context(), r)
			w.WriteHeader(status)
			fmt.Fprint(w, message)
		})).ServeHTTP(w, r)
	})

	// Start the server (this is not typically used in cloud functions, but local development)
	fmt.Println("Starting server...")
	http.ListenAndServe(":8080", nil)
}
