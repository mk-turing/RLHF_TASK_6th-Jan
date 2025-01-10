package main

import (
	"fmt"
	"sort"
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
	Name       string
	Attendees  map[string][]Attendee // Maps categories to attendees
	Categories []string               // Stores categories in sorted order
	lock       sync.Mutex
}

func (e *Event) RegisterAttendee(attendee Attendee) {
	e.lock.Lock()
	defer e.lock.Unlock()
	// Add attendee to the corresponding category
	e.Attendees[attendee.Category] = append(e.Attendees[attendee.Category], attendee)

	// Add new category to the sorted slice if not present
	if index := sort.SearchStrings(e.Categories, attendee.Category); index == len(e.Categories) {
		e.Categories = append(e.Categories, attendee.Category)
	}
}

func (e *Event) GetAttendees() map[string][]Attendee {
	e.lock.Lock()
	defer e.lock.Unlock()
	attendees := make(map[string][]Attendee, len(e.Attendees))
	for category, attendeesList := range e.Attendees {
		attendees[category] = append(attendees[category], attendeesList...)
	}
	return attendees
}

func (e *Event) CountAttendeesByCategory(category string) int {
	return len(e.Attendees[category])
}

func (e *Event) SortAttendeesByCategory() map[string][]Attendee {
	e.lock.Lock()
	defer e.lock.Unlock()

	sortedAttendees := make(map[string][]Attendee, len(e.Attendees))
	for _, category := range e.Categories {
		sortedAttendees[category] = append(sortedAttendees[category], e.Attendees[category]...)
		sort.Slice(sortedAttendees[category], func(i, j int) bool {
			// Sort attendees first by name and then by email
			a, b := sortedAttendees[category][i], sortedAttendees[category][j]
			return a.Name < b.Name || (a.Name == b.Name && a.Email < b.Email)
		})
	}
	return sortedAttendees
}

func (e *Event) GetAttendeeByEmail(email string) *Attendee {
	for _, attendees := range e.Attendees {
		for _, attendee := range attendees {
			if attendee.Email == email {
				return &attendee
			}
		}
	}
	return nil
}

func main() {
	event1 := Event{Name: "Tech Conference"}

	event1.RegisterAttendee(Attendee{Name: "Alice", Email: "alice@example.com", Age: 28, Gender: "F", Category: "VIP"})
	event1.RegisterAttendee(Attendee{Name: "Bob", Email: "bob@example.com", Age: 30, Gender: "M", Category: "Regular"})
	event1.RegisterAttendee(Attendee{Name: "Charlie", Email: "charlie@example.com", Age: 25, Gender: "M", Category: "VIP"})
	event1.RegisterAttendee(Attendee{Name: "David", Email: "david@example.com", Age: 35, Gender: "M", Category: "VIP"})
	event1.RegisterAttendee(Attendee{Name: "Emma", Email: "emma@example.com", Age: 27, Gender: "F", Category: "Regular"})

	fmt.Println("--- Event Attendees ---")
	fmt.Println(event1.GetAttendees())

	fmt.Println("\n--- Attendees by Category (sorted) ---")
	sortedAttendees := event1.SortAttendeesByCategory()
	for category, attendees := range sortedAttendees {
		fmt.Printf("%s:\n", category)