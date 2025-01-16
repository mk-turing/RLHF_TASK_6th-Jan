package main

import (
	"fmt"
	"net/http"
	"strings"
)

// Middleware represents a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// RouteConfig contains the path and middleware for a single route.
type RouteConfig struct {
	Path        string
	Middlewares []Middleware
}

// Gateway holds maps of middleware components and route configurations.
type Gateway struct {
	Middlewares map[string]Middleware
	Routes      map[string]RouteConfig
}

// NewGateway creates a new Gateway instance with empty maps.
func NewGateway() *Gateway {
	return &Gateway{
		Middlewares: make(map[string]Middleware),
		Routes:      make(map[string]RouteConfig),
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

// AddRouteConfig adds or updates a route configuration in the gateway.
func (g *Gateway) AddRouteConfig(path string, middlewareNames ...string) {
	config := RouteConfig{
		Path:        path,
		Middlewares: []Middleware{},
	}

	for _, name := range middlewareNames {
		if middleware, ok := g.GetMiddleware(name); ok {
			config.Middlewares = append(config.Middlewares, middleware)
		}
	}

	g.Routes[path] = config
}

// GetRouteConfig retrieves a route configuration from the gateway based on the path.
func (g *Gateway) GetRouteConfig(path string) (RouteConfig, bool) {
	config, ok := g.Routes[path]
	return config, ok
}

// HandleRequest processes a request using the appropriate middleware stack based on the route.
func (g *Gateway) HandleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if config, ok := g.GetRouteConfig(path); ok {
		var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, World!")
		})

		for _, middleware := range config.Middlewares {
			handler = middleware(handler)
		}

		handler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// Example Middleware Functions:
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Logging request: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Perform authentication logic here
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func sampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Inside sample middleware.\n")
		next.ServeHTTP(w, r)
	})
}

func main() {
	gateway := NewGateway()