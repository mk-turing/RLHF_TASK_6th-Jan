package main

import (
	"fmt"
	"net/http"
)

// Middleware represents a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// RoutingPolicy defines a set of conditions and middleware for a given route.
type RoutingPolicy struct {
	Condition  func(*http.Request) bool
	Middleware []string
}

// Gateway is a struct that holds a map of routing policies and middleware functions.
type Gateway struct {
	RoutingPolicies map[string]RoutingPolicy
	Middlewares     map[string]Middleware
}

// NewGateway creates a new Gateway instance with an empty map of routing policies.
func NewGateway() *Gateway {
	return &Gateway{
		RoutingPolicies: make(map[string]RoutingPolicy),
		Middlewares:     make(map[string]Middleware),
	}
}

// AddRoutingPolicy adds or updates a routing policy for the gateway.
func (g *Gateway) AddRoutingPolicy(path string, policy RoutingPolicy) {
	g.RoutingPolicies[path] = policy
}

// GetRoutingPolicy retrieves a routing policy from the gateway by path.
func (g *Gateway) GetRoutingPolicy(path string) (RoutingPolicy, bool) {
	policy, ok := g.RoutingPolicies[path]
	return policy, ok
}

// AddMiddleware adds a middleware function to the gateway.
func (g *Gateway) AddMiddleware(name string, middleware Middleware) {
	g.Middlewares[name] = middleware
}

// GetMiddleware retrieves a middleware function by name.
func (g *Gateway) GetMiddleware(name string) (Middleware, bool) {
	middleware, ok := g.Middlewares[name]
	return middleware, ok
}

// HandleRequest processes a request using the specified middleware components for the matching route.
func (g *Gateway) HandleRequest(w http.ResponseWriter, r *http.Request) {
	policy, ok := g.GetRoutingPolicy(r.URL.Path)
	if !ok {
		http.Error(w, "No matching route found", http.StatusNotFound)
		return
	}

	if !policy.Condition(r) {
		http.Error(w, "Request does not match condition", http.StatusForbidden)
		return
	}

	// Initialize the handler as an http.HandlerFunc
	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Handling path: %s", r.URL.Path)
	})

	// Apply the middleware to the handler iteratively
	for _, name := range policy.Middleware {
		middleware, ok := g.GetMiddleware(name)
		if !ok {
			http.Error(w, fmt.Sprintf("Middleware '%s' not found", name), http.StatusInternalServerError)
			return
		}
		// Convert handler to http.Handler and apply middleware
		handler = middleware(handler) // handler is already of type http.Handler, so no cast needed
	}

	// Serve the request using the final handler
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

// Example Condition Functions
func isTokenPresent(r *http.Request) bool {
	_, present := r.Header["Authorization"]
	return present
}

func hasQueryParam(r *http.Request, param string) bool {
	return r.URL.Query().Get(param) != ""
}

func main() {
	gateway := NewGateway()

	// Add or update middleware components dynamically
	gateway.AddMiddleware("logger", loggerMiddleware)
	gateway.AddMiddleware("auth", authMiddleware)

	// Define routing policies with conditions
	gateway.AddRoutingPolicy("/api/users", RoutingPolicy{
		Condition:  isTokenPresent,
		Middleware: []string{"logger", "auth"},
	})

	gateway.AddRoutingPolicy("/api/public", RoutingPolicy{
		Condition: func(r *http.Request) bool {
			return hasQueryParam(r, "param1")
		},
		Middleware: []string{"logger"},
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gateway.HandleRequest(w, r)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
