package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type User struct {
	Username string
	Password string
	Token    string
	Role     string // New field to add roles (e.g., "admin", "voter")
}

type Poll struct {
	ID      int
	Question string
	Answers []string
	Created time.Time
}

var users []User
var sessions []string
var polls []Poll

func main() {
	users = append(users, User{Username: "admin", Password: "password", Role: "admin"})
	// Add other users, such as voters

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/poll", pollHandler)
	http.HandleFunc("/create_poll", createPollHandler)
	http.HandleFunc("/delete_poll/{poll_id}", deletePollHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	for _, user := range users {
		if user.Username == username && user.Password == password {
			token := generateToken()
			user.Token = token
			sessions = append(sessions, token)
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   token,
				Expires: time.Now().Add(30 * time.Minute),
			})
			if user.Role == "admin" {
				http.Redirect(w, r, "/create_poll", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/poll", http.StatusSeeOther)
			}
			return
		}
	}

	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
}

func pollHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := cookie.Value

	if isValidToken(token) {
		jsonPollList, err := json.Marshal(polls)
		if err != nil {
			http.Error(w, "Error marshaling polls", http.StatusInternalServerError)
			return
		}
		w.Write(jsonPollList)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func createPollHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := cookie.Value

	if isValidToken(token) {
		var poll Poll
		if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		polls = append(polls, poll)
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func deletePollHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := cookie.Value
	vars := r.URL.Query()
	pollIDStr := vars.Get("poll_id")
	pollID, err := strconv.Atoi(pollIDStr)
	if err != nil {
		http.Error(w, "Invalid poll ID", http.StatusBadRequest)
		return
	}

	if isValidToken(token) && strings.Contains(userRoles(token), "admin") {
		var index int
		for i, poll := range polls {
			if poll.ID == pollID {
				index = i
				break
			}
		}
		if index != -1 {
			polls = append(polls[:index], polls[index+1:]...)
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Poll not found", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func isValidToken(token string) bool {
	for _, session := range sessions {
		if session == token {
			return true
		}
	}
	return false
}

func userRoles(token string) string {
	for _, user := range users {
		if user.Token == token {
			return user.Role
		}
	}
	return ""
}

func generateToken() string {
	return "random-token-string" // Use a CRNG in production
}

func strconv.Atoi(pollIDStr string) (int, error) {
	return 0, nil // Implement proper error handling here
}