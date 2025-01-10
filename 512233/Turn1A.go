package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Attendee represents an attendee of the event
type Attendee struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Attending bool   `json:"attending"`
}

// RegisterAttendee adds an attendee to the registration list
func RegisterAttendee(attendees *[]Attendee, name, email, phone string) {
	*attendees = append(*attendees, Attendee{Name: name, Email: email, Phone: phone, Attending: true})
}

// GenerateReport generates a report of registered attendees
func GenerateReport(attendees []Attendee) string {
	report := "Attendees:\n"
	for _, attendee := range attendees {
		report += fmt.Sprintf("Name: %s, Email: %s, Phone: %s\n", attendee.Name, attendee.Email, attendee.Phone)
	}
	return report
}

// ExtractInsights extracts detailed insights from attendee data
func ExtractInsights(attendees []Attendee) (emails []string, phoneNumbers []string) {
	for _, attendee := range attendees {
		emails = append(emails, attendee.Email)
		phoneNumbers = append(phoneNumbers, attendee.Phone)
	}
	return
}

// Main function to test the event registration application
func main() {
	attendees := make([]Attendee, 0)

	RegisterAttendee(&attendees, "Alice Smith", "alice@example.com", "123-456-7890")
	RegisterAttendee(&attendees, "Bob Johnson", "bob@example.com", "098-765-4321")

	fmt.Println("Event Registration Report:")
	fmt.Println(GenerateReport(attendees))

	emails, phoneNumbers := ExtractInsights(attendees)
	fmt.Println("\nInsights:")
	fmt.Println("Emails:", emails)
	fmt.Println("Phone Numbers:", phoneNumbers)

	// Serialize attendees for storage or transmission
	data, err := json.Marshal(attendees)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nSerialized Attendee Data:", string(data))
}
