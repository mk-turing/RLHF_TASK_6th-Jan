package main

import (
	"fmt"
	"sync"
)

type Attendee struct {
	Name   string
	Email  string
	Age    int
	Gender string
}
type Event struct {
	Name      string
	Attendees []Attendee
	lock      sync.Mutex
}

func (e *Event) RegisterAttendee(attendee Attendee) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.Attendees = append(e.Attendees, attendee)
}

func (e *Event) GetAttendees() []Attendee {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.Attendees
}

func (e *Event) CountAttendees() int {
	return len(e.GetAttendees())
}

func (e *Event) GetAttendeeByEmail(email string) *Attendee {
	for _, attendee := range e.GetAttendees() {
		if attendee.Email == email {
			return &attendee
		}
	}
	return nil
}

func main() {
	event := Event{Name: "Golang Workshop"}
	event.RegisterAttendee(Attendee{Name: "Alice", Email: "alice@example.com", Age: 28, Gender: "F"})
	event.RegisterAttendee(Attendee{Name: "Bob", Email: "bob@example.com", Age: 30, Gender: "M"})
	event.RegisterAttendee(Attendee{Name: "Charlie", Email: "charlie@example.com", Age: 25, Gender: "M"})
	// Register more attendees...

	fmt.Println("Total Attendees:", event.CountAttendees())
	attendee := event.GetAttendeeByEmail("bob@example.com")
	if attendee != nil {
		fmt.Println("Attendee Found:", attendee.Name)
	} else {
		fmt.Println("Attendee not found.")
	}

	// Data analysis using slices
	attendeesByGender := make(map[string]int)
	for _, attendee := range event.GetAttendees() {
		attendeesByGender[attendee.Gender]++
	}
	fmt.Println("Attendees by Gender:")
	for gender, count := range attendeesByGender {
		fmt.Printf("%s: %d\n", gender, count)
	}
}
