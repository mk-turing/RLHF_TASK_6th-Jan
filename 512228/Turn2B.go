package main

import (
	"fmt"
	"net/http"
)

type User struct {
	Username string
	Password string
	Token    string
	Role     string // added Role field
}

var users []User
var sessions []string

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
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Same logic as before
}

func pollHandler(w http.ResponseWriter, r *http.Request) {
	// Validate user session and role using getUserFromToken function
	user := getUserFromToken(r)
	if user == nil || user.Role != "voter" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Authenticated voter, display poll page
	fmt.Fprintf(w, "Welcome to the poll page, %s!", user.Username)
}

func createPollHandler(w http.ResponseWriter, r *http.Request) {
	// Validate user session and role using getUserFromToken function
	user := getUserFromToken(r)
	if user == nil || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Authenticated admin, handle poll creation logic
	fmt.Fprintf(w, "Welcome, %s! You can create a new poll here.", user.Username)
}

func deletePollHandler(w http.ResponseWriter, r *http.Request) {
	// Validate user session and role using getUserFromToken function
	user := getUserFromToken(r)
	if user == nil || user.Role != "admin" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Authenticated admin, handle poll deletion logic
	fmt.Fprintf(w, "Welcome, %s! You can delete a poll here.", user.Username)
}

func getUserFromToken(r *http.Request) *User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	token := cookie.Value

	// Validate the token using the sessions slice and retrieve the corresponding user
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
