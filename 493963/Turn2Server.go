package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func yourEndpointHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Respond with a simple message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	// Create a new HTTP server
	http.HandleFunc("/your-endpoint", yourEndpointHandler)

	// Define the server's address and port
	serverAddress := ":8080"
	log.Printf("Starting server at %s...\n", serverAddress)

	// Start the server
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
