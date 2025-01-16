package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Middleware represents a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// RoutingPolicy defines a set of conditions and middleware for a given route.
type RoutingPolicy struct {
	Condition    func(*http.Request) bool
	Middleware   []string
	RateLimit    *RateLimitConfig
	Authenticate bool
}

// RateLimitConfig represents a configuration for rate-limiting a specific route.
type RateLimitConfig struct {
	MaxRequests  int
	WindowPeriod time.Duration
}

// RateLimit is a struct that manages request rates.
type RateLimit struct {
	lastRequest time.Time
	mu          sync.Mutex
	maxRequests int
	window      time.Duration
	count       int
}

// NewRateLimit creates a new RateLimit instance with the specified configuration.
func NewRateLimit(maxRequests int, window time.Duration) *RateLimit {
	return &RateLimit{
		lastRequest: time.Now(),
		mu:          sync.Mutex{},
		maxRequests: maxRequests,
		window:      window,
		count:       0,
	}
}

// IsAllowed checks if a request can be processed within the rate limit.
func (r *RateLimit) IsAllowed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if time.Now().Sub(r.lastRequest) >= r.window {
		r.count = 0
	}

	if r.count >= r.maxRequests {
		return false
	}

	r.count++
	r.lastRequest = time.Now()
	return true
}

// Gateway is a struct that holds a map of routing policies and middleware functions.
type Gateway struct {
	RoutingPolicies map[string]RoutingPolicy
	Middlewares     map[string]Middleware
	RateLimits      map[string]*RateLimit
}

// NewGateway creates a new Gateway instance with an empty map of routing policies.
func NewGateway() *Gateway {
	return &Gateway{
		RoutingPolicies: make(map[string]RoutingPolicy),
		Middlewares:     make(map[string]Middleware),
		RateLimits:      make(map[string]*RateLimit),
	}
}

// AddRoutingPolicy adds or updates a routing policy for the gateway.
func (g *Gateway) AddRoutingPolicy(path string, policy RoutingPolicy) {
	g.RoutingPolicies[path] = policy
	if policy.RateLimit != nil {
		g.RateLimits[path] = NewRateLimit(policy.RateLimit.MaxRequests, policy.RateLimit.WindowPeriod)
	}
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

	// Check authentication if required
	if policy.Authenticate && !isTokenValid(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check rate limit
	if policy.RateLimit != nil && !g.RateLimits[r.URL.Path].IsAllowed() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
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
		handler = middleware(handler)
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

func isTokenValid(r *http.Request) bool {
	_, present := r.Header["Authorization"]
	return present
}

func main() {
	gateway := NewGateway()

	// Add or update middleware components dynamically
	gateway.AddMiddleware("logger", loggerMiddleware)

	// Define routing policies with conditions, authentication, and rate limits
	gateway.AddRoutingPolicy("/api/users", RoutingPolicy{
		Condition:    func(r *http.Request) bool { return true },
		Middleware:   []string{"logger"},
		RateLimit:    &RateLimitConfig{MaxRequests: 5, WindowPeriod: time.Second * 10},
		Authenticate: true,
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gateway.HandleRequest(w, r)
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
