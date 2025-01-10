package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
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
	Attendees map[string][]Attendee `json:"attendees"`
	lock      sync.Mutex
}

var events = make(map[string]Event)
var db *sql.DB

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./events.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS attendees (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_category TEXT,
		name TEXT,
		email TEXT,
		age INTEGER,
		gender TEXT,
		registered_at TEXT
	);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (e *Event) RegisterAttendee(attendee Attendee) {
	e.lock.Lock()
	defer e.lock.Unlock()

	attendee.ID = len(e.Attendees[attendee.Category]) + 1
	attendee.RegisteredAt = time.Now()
	e.Attendees[attendee.Category] = append(e.Attendees[attendee.Category], attendee)

	// Save to database
	saveAttendeeToDB(attendee)
}

func saveAttendeeToDB(attendee Attendee) {
	_, err := db.Exec(`INSERT OR REPLACE INTO attendees (event_category, name, email, age, gender, registered_at) VALUES (?, ?, ?, ?, ?, ?)`,
		attendee.Category, attendee.Name, attendee.Email, attendee.Age, attendee.Gender, attendee.RegisteredAt.Format(time.RFC3339))
	if err != nil {
		fmt.Println("Error saving attendee:", err)
	}
}

func loadAttendeesFromDB(category string) []Attendee {
	rows, err := db.Query(`SELECT * FROM attendees WHERE event_category = ? ORDER BY id ASC`, category)
	if err != nil {
		fmt.Println("Error loading attendees:", err)
		return nil
	}
	defer rows.Close()
	var attendees []Attendee
	for rows.Next() {
		var attendee Attendee
		err := rows.Scan(&attendee.ID, &attendee.Category, &attendee.Name, &attendee.Email, &attendee.Age, &attendee.Gender, &attendee.RegisteredAt)
		if err != nil {
			fmt.Println("Error scanning attendee:", err)
			continue
		}
		attendees = append(attendees, attendee)
	}
	return attendees
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

func (e *Event) GetAttendees() map[string][]Attendee {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.Attendees
}

func getEventAttendeesHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	event, ok := events[category]
	if !ok {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	attendees := loadAttendeesFromDB(category)
	event.Attendees[category] = attendees

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event.GetAttendees())
}

func getAttendanceStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		Total   int `json:"total"`
		Vip     int `json:"vip"`
		Regular int `json:"regular"`
	}

	rows, err := db.Query(`SELECT COUNT(*) as total_count, event_category, COUNT(DISTINCT id) as attendee_count FROM attendees GROUP BY event_category`)
	if err != nil {
		fmt.Println("Error querying statistics:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		var category string
		if err := rows.Scan(&count, &category, &count); err != nil {
			fmt.Println("Error scanning statistics:", err)
			continue
		}
		switch category {
		case "VIP":
			stats.Vip += count
		case "Regular":
			stats.Regular += count
		default:
			stats.Total += count
		}
	}

	stats.Total += stats.Vip + stats.Regular

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func main() {
	var err error
	if db, err = connectDB(); err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/register", registerAttendeeHandler)
	http.HandleFunc("/attendees", getEventAttendeesHandler)
	http.HandleFunc("/statistics", getAttendanceStatisticsHandler)

	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
