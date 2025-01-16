package main

import (
	"fmt"
	"net/http"
)

// Middleware interface to define a middleware function
type Middleware func(http.Handler) http.Handler

// MiddlewareMap to store registered middleware functions
var MiddlewareMap = map[string]Middleware{}

// RegisterMiddleware function to add a middleware to the map
func RegisterMiddleware(name string, mw Middleware) {
	MiddlewareMap[name] = mw
	fmt.Printf("Registered middleware: %v\n", name)
}

// ValidateRequest middleware to check if the request has a specific header
func ValidateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Valid-Key") == "" {
			http.Error(w, "Missing X-Valid-Key", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// TransformRequest middleware to modify the request path
func TransformRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/transformed/" + r.URL.Path
		next.ServeHTTP(w, r)
	})
}

// APIGatewayHandler to process requests using registered middleware
func APIGatewayHandler(w http.ResponseWriter, r *http.Request) {
	var handlers []http.Handler

	// Example of finding middleware based on request properties
	for name, mw := range MiddlewareMap {
		// Deciding the logic for applying middleware based on criteria like URL path
		if r.URL.Path == "/api/example" {
			handlers = append(handlers, mw)
		}
	}

	// Apply middleware in the correct order
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Original path: %v Transformed path: %v\n", r.URL.Path, r.URL.Path)
	})

	for i := len(handlers) - 1; i >= 0; i-- {
		next = handlers[i](next)
	}

	next.ServeHTTP(w, r)
}

func main() {
	// Register middleware functions
	RegisterMiddleware("validateRequest", ValidateRequest)
	RegisterMiddleware("transformRequest", TransformRequest)

	// Define routes and handlers
	http.HandleFunc("/api/example", APIGatewayHandler)
	http.ListenAndServe(":8080", nil)
}
