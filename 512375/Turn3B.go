package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
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
	rateLimiter     *RateLimiter
	authHandler    http.Handler
}

// NewGateway creates a new Gateway instance with an empty map of routing policies and initializes rate limiter.
func NewGateway() *Gateway {
	return &Gateway{
		RoutingPolicies: make(map[string]RoutingPolicy),
		Middlewares:     make(map[string]Middleware),
		rateLimiter:     NewRateLimiter(),
		authHandler:    http.HandlerFunc(defaultAuthHandler), // Provide a default auth handler or use a real implementation
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

// SetAuthHandler sets the custom authentication handler.
func (g *Gateway) SetAuthHandler(authHandler http.Handler) {
	g.authHandler = authHandler
}

// HandleRequest processes a request using the specified middleware components for the matching route.
func (g *Gateway) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// ... (Same code as before)

	// Apply authentication middleware
	g.authHandler.ServeHTTP(w, r)
	if r.Context().Err() != nil {
		return
	}

	// Apply rate-limiting middleware
	if !g.rateLimiter.Allow(r) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// ... (Same code as before)
}

// Example Authentication Middleware
func defaultAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Basic authentication for demonstration purposes
	username, password, ok := r.BasicAuth()
	if !ok || username != "user" || password != "password" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Set any necessary request context values for subsequent middleware
	r = r.WithContext(context.WithValue(r.Context(), "user", username))
}

// RateLimiter is a simple rate limiter for demonstration purposes.
// In a real application, you should use a distributed rate limiter like Redis or Consul.
type RateLimiter struct {
	mu sync.Mutex
	// Map of user identifiers (e.g., API keys) to their respective rate limiters
	limits map[string]*Rate
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limits: make(map[string]*Rate),
	}
}

// Allow checks if a request is allowed based on the rate limit for the user.
func (rl *RateLimiter) Allow(r *http.Request) bool {
	// Extract the user identifier from the request (e.g., API key)