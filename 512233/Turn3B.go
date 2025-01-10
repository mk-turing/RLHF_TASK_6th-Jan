package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
	"time"
)

type Attendee struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Gender   string `json:"gender"`
	Category string `json:"category"`
}

type Event struct {
	Name      string
	Attendees map[string][]Attendee
	Analytics map[time.Time]int // Maps registration time to attendee count
	lock      sync.Mutex
}

var event Event

func init() {
	event = Event{
		Name:      "Golang Workshop",
		Attendees: make(map[string][]Attendee),
		Analytics: make(map[time.Time]int),
	}
}

func (e *Event) RegisterAttendee(attendee Attendee) {
	e.lock.Lock()
	defer e.lock.Unlock()

	// Register the attendee
	e.Attendees[attendee.Category] = append(e.Attendees[attendee.Category], attendee)

	// Update real-time analytics
	currentTime := time.Now().Truncate(time.Minute)
	e.Analytics[currentTime]++
}

func registerAttendeeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var attendee Attendee
	if err := json.NewDecoder(r.Body).Decode(&attendee); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event.RegisterAttendee(attendee)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Attendee registered successfully"})
}

func attendeeStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	event.lock.Lock()
	defer event.lock.Unlock()

	var statistics []map[string]interface{}

	for time, count := range event.Analytics {
		statistics = append(statistics, map[string]interface{}{
			"time":  time.Format("2006-01-02 15:04"),
			"count": count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statistics)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/attendee", registerAttendeeHandler).Methods("POST")
	r.HandleFunc("/statistics", attendeeStatisticsHandler).Methods("GET")

	fmt.Println("Event Registration API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
