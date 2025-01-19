package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Middleware func(http.Handler) http.Handler

// AddTimeout adds a timeout to the request context
func AddTimeout(duration time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Entry-point for the cloud function
func EntryPoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Optimized Cloud Function!"))
}

func main() {
	r := mux.NewRouter()
	middlewares := []Middleware{
		AddTimeout(2 * time.Second), // Apply a timeout of 2 seconds
	}

	handler := http.Handler(http.HandlerFunc(EntryPoint))
	for _, m := range middlewares {
		handler = m(handler)
	}

	r.Handle("/", handler)

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
