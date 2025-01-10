package main

import (
	"fmt"
	"net/http"
)

type User struct {
	Username string
	Password string
	Token    string
	Role     string
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
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/poll", pollHandler)
	http.HandleFunc("/createPoll", createPollHandler)
	http.HandleFunc("/deletePoll", deletePollHandler)
	http.HandleFunc("/vote", voteHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Same logic as before
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
					return &user
				}
			}
		}
	}
	return nil
}

// Helper function to check if a slice contains a specific element
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
