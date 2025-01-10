package main

import (
	"fmt"
	"sync"
)

type Attendee struct {
	Name     string
	Email    string
	Age      int
	Gender   string
	Category string
}

type Event struct {
	Name      string
	Attendees map[string][]Attendee // Map to store attendees by category
	lock      sync.Mutex
}

func (e *Event) RegisterAttendee(attendee Attendee) {
	e.lock.Lock()
	defer e.lock.Unlock()

	// Check if the category already exists, if not, create it
	if _, ok := e.Attendees[attendee.Category]; !ok {
		e.Attendees[attendee.Category] = make([]Attendee, 0)
	}

	// Append the attendee to the appropriate category
	e.Attendees[attendee.Category] = append(e.Attendees[attendee.Category], attendee)
}

func (e *Event) GetAttendees() map[string][]Attendee {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.Attendees
}

func (e *Event) CountAttendees() int {
	total := 0
	for _, categoryAttendees := range e.GetAttendees() {
		total += len(categoryAttendees)
	}
	return total
}

func (e *Event) SortAttendeesByPriority() {
	for _, categoryAttendees := range e.GetAttendees() {
		// Implement custom sorting based on priority
		// For example, sort by category and then by name
		sortedAttendees := make([]Attendee, len(categoryAttendees))
		copy(sortedAttendees, categoryAttendees)

		// Sorting logic
		for i := 0; i < len(sortedAttendees)-1; i++ {
			for j := 0; j < len(sortedAttendees)-i-1; j++ {
				// Compare by category first, then by name
				if sortedAttendees[j].Category > sortedAttendees[j+1].Category {
					sortedAttendees[j], sortedAttendees[j+1] = sortedAttendees[j+1], sortedAttendees[j]
				} else if sortedAttendees[j].Category == sortedAttendees[j+1].Category && sortedAttendees[j].Name > sortedAttendees[j+1].Name {
					sortedAttendees[j], sortedAttendees[j+1] = sortedAttendees[j+1], sortedAttendees[j]
				}
			}
		}

		// Assign sorted attendees back to the category
		e.Attendees[categoryAttendees[0].Category] = sortedAttendees
	}
}

func (e *Event) ReportAttendees() {
	e.SortAttendeesByPriority()
	fmt.Println("Attendees for Event:", e.Name)
	fmt.Println("--------------------")

	for category, attendees := range e.GetAttendees() {
		fmt.Printf("Category: %s\n", category)
		for _, attendee := range attendees {
			fmt.Printf("- %s (%s)\n", attendee.Name, attendee.Email)
		}
		fmt.Println()
	}
}

func main() {
	event := Event{
		Name:      "Golang Workshop",
		Attendees: make(map[string][]Attendee), // Initialize the map here
	}

	// Register attendees
	event.RegisterAttendee(Attendee{Name: "Alice", Email: "alice@example.com", Age: 28, Gender: "F", Category: "VIP"})
	event.RegisterAttendee(Attendee{Name: "Bob", Email: "bob@example.com", Age: 30, Gender: "M", Category: "Regular"})
	event.RegisterAttendee(Attendee{Name: "Charlie", Email: "charlie@example.com", Age: 25, Gender: "M", Category: "VIP"})
	event.RegisterAttendee(Attendee{Name: "Diana", Email: "diana@example.com", Age: 23, Gender: "F", Category: "Regular"})

	fmt.Println("Total Attendees:", event.CountAttendees())

	event.ReportAttendees()
}
