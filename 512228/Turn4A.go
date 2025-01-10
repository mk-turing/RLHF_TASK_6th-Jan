package main

import (
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Username string
	Password string
	Token    string
	Role     string
	LastUsed time.Time // Track the last time the session was used
}

var users []User
var sessions []string
var votes = make(map[string][]string) // Use a map to track votes by poll ID

func main() {
	users = append(users,
		User{Username: "admin", Password: "password", Role: "admin"},
		User{Username: "voter1", Password: "voter1", Role: "voter"},
		User{Username: "voter2", Password: "voter2", Role: "voter"},
	)

	go manageSessionTimeouts() // Start a goroutine to manage session timeouts

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/poll", pollHandler)
	http.HandleFunc("/createPoll", createPollHandler)
	http.HandleFunc("/deletePoll", deletePollHandler)
	http.HandleFunc("/vote", voteHandler)
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

	for i, user := range users {
		if user.Username == username && user.Password == password {
			// Generate a new token and update user's last used time
			token := generateToken()
			users[i].Token = token
			user.LastUsed = time.Now()
			sessions = append(sessions, token)
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    token,
				Expires:  time.Now().Add(30 * time.Minute), // Set cookie expiration time
				HttpOnly: true,                             // Prevent client-side access to the cookie
			})
			http.Redirect(w, r, "/poll", http.StatusSeeOther)
			return
		}
	}
	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
}

func pollHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromToken(r)
	if user == nil || user.Role != "voter" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Authenticated voter, display poll page
	fmt.Fprintf(w, "Welcome to the poll page, %s!", user.Username)
}

func createPollHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromToken(r)
	if user == nil || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Authenticated admin, handle poll creation logic
	fmt.Fprintf(w, "Welcome, %s! You can create a new poll here.", user.Username)
}

func deletePollHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromToken(r)
	if user == nil || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Authenticated admin, handle poll deletion logic
	fmt.Fprintf(w, "Welcome, %s! You can delete a poll here.", user.Username)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromToken(r)
	if user == nil || user.Role != "voter" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pollID := r.FormValue("poll_id")
	option := r.FormValue("option")

	if pollID == "" || option == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if user has already voted on this poll
	if _, ok := votes[pollID]; ok {
		if contains(votes[pollID], user.Username) {
			http.Error(w, "You have already voted for this poll", http.StatusForbidden)
			return
		}
	}

	// Record the vote
	votes[pollID] = append(votes[pollID], user.Username)
	fmt.Fprintf(w, "Thank you for voting, %s! Your vote for '%s' has been recorded.", user.Username, option)
}

func getUserFromToken(r *http.Request) *User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	token := cookie.Value

	for _, session := range sessions {
		if session == token {
			for _, user := range users {
				if user.Token == token {
					// Update last used time if session is valid
					user.LastUsed = time.Now()
					return &user
				}
			}
		}
	}
	return nil
}

func manageSessionTimeouts() {
	sessionTimeout := 30 * time.Minute // Define session timeout duration
	for {
		time.Sleep(sessionTimeout) // Wait for the timeout duration
		inactiveSessions := []string{}
		for _, session := range sessions {
			for _, user := range users {
				if user.Token == session {
					if time.Since(user.LastUsed) > sessionTimeout {
						inactiveSessions = append(inactiveSessions, session)
						break
					}
				}
			}
		}
		for _, session := range inactiveSessions {
			// Remove the inactive session from the slices
			var newSessions []string
			var newUsers []User
			for _, s := range sessions {
				if s != session {
					newSessions = append(newSessions, s)
				}
			}
			for _, u := range users {
				if u.Token != session {
					newUsers = append(newUsers, u)
				}
			}
			sessions = newSessions
			users = newUsers
			fmt.Println("Session", session, "invalidated due to inactivity.")
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func generateToken() string {
	// Generate a unique token using a cryptographically secure random number generator (CRNG)
	// For simplicity, we'll use a random string in this example
	return "random-token-string"
}
