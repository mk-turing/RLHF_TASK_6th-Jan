package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Dynamic scaling function to simulate worker scaling
func scaleWorkers() {
	// Placeholder function to scale up or down workers based on workload.
	// In a real implementation, this would interact with cloud services (e.g., AWS, GCP)
	// to adjust the number of workers dynamically.
	fmt.Println("Scaling workers dynamically based on load...")
}

// Middleware to simulate request processing
func requestProcessingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate varying resource demands of middleware
		if r.URL.Path == "/heavy-task" {
			// Trigger worker scaling for heavy tasks
			scaleWorkers()
		}
		next.ServeHTTP(w, r)
	})
}

// Example cloud function
func cloudFunction(w http.ResponseWriter, r *http.Request) {
	// Simulate response generation
	time.Sleep(time.Millisecond * 200) // Simulate task processing
	w.Write([]byte("Task Processed"))
}

func main() {
	r := mux.NewRouter()

	// Apply dynamic scaling middleware
	r.Use(requestProcessingMiddleware)

	// Define routes
	r.HandleFunc("/", cloudFunction)
	r.HandleFunc("/heavy-task", cloudFunction)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
