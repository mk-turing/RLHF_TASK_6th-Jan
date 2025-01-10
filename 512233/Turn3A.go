package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Attendee struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Age          int       `json:"age"`
	Gender       string    `json:"gender"`
	Category     string    `json:"category"`
	RegisteredAt time.Time `json:"registeredAt"`
}

type Event struct {
	Name      string
	Attendees []Attendee `json:"attendees"`
	lock      sync.Mutex
}

var events = make(map[string]Event)

func (e *Event) RegisterAttendee(attendee Attendee) {
	e.lock.Lock()
	defer e.lock.Unlock()
	attendee.ID = len(e.Attendees) + 1
	attendee.RegisteredAt = time.Now()
	e.Attendees = append(e.Attendees, attendee)
}

func (e *Event) GetAttendees() []Attendee {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.Attendees
}

func (e *Event) CountAttendees() int {
	return len(e.Attendees)
}

func registerAttendeeHandler(w http.ResponseWriter, r *http.Request) {
	var attendee Attendee
	if err := json.NewDecoder(r.Body).Decode(&attendee); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	event, ok := events[attendee.Category]
	if !ok {
		event = Event{Name: attendee.Category}
		events[attendee.Category] = event
	}

	event.RegisterAttendee(attendee)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attendee)
}

func getEventAttendeesHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	event, ok := events[category]
	if !ok {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event.GetAttendees())
}

func getAttendanceStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		Total   int `json:"total"`
		Vip     int `json:"vip"`
		Regular int `json:"regular"`
	}

	for _, event := range events {
		stats.Total += event.CountAttendees()
		for _, attendee := range event.Attendees {
			if attendee.Category == "VIP" {
				stats.Vip++
			} else {
				stats.Regular++
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func main() {
	http.HandleFunc("/register", registerAttendeeHandler)
	http.HandleFunc("/attendees", getEventAttendeesHandler)
	http.HandleFunc("/statistics", getAttendanceStatisticsHandler)

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
