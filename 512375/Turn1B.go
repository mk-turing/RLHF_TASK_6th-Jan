package main

import (
	"fmt"
	"net/http"
)

// Middleware represents a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// Gateway is a struct that holds a map of middleware components.
type Gateway struct {
	Middlewares map[string]Middleware
}

// NewGateway creates a new Gateway instance with an empty map of middleware.
func NewGateway() *Gateway {
	return &Gateway{
		Middlewares: make(map[string]Middleware),
	}
}

// AddMiddleware adds or updates a middleware component in the gateway.
func (g *Gateway) AddMiddleware(name string, middleware Middleware) {
	g.Middlewares[name] = middleware
}

// GetMiddleware retrieves a middleware component from the gateway by name.
func (g *Gateway) GetMiddleware(name string) (Middleware, bool) {
	middleware, ok := g.Middlewares[name]
	return middleware, ok
}

// HandleRequest processes a request using the specified middleware components.
func (g *Gateway) HandleRequest(w http.ResponseWriter, r *http.Request, middlewareNames []string) {
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	for _, name := range middlewareNames {
		middleware, ok := g.GetMiddleware(name)
		if !ok {
			fmt.Fprintf(w, "Middleware '%s' not found.", name)
			return
		}
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// Example Middleware Functions
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Logging request: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Perform authentication logic here
		// For simplicity, let's just assume authentication is successful
		next.ServeHTTP(w, r)
	})
}

func main() {
	gateway := NewGateway()

	// Add or update middleware components dynamically
	gateway.AddMiddleware("logger", loggerMiddleware)
	gateway.AddMiddleware("auth", authMiddleware)

	// Define the middleware stack for a specific route
	middlewareStack := []string{"logger", "auth"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gateway.HandleRequest(w, r, middlewareStack)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
