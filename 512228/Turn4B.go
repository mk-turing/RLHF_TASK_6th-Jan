package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	Username   string
	Password   string
	Token      string
	Role       string
	LastActive time.Time
}

var users []User
var sessionExpirationDuration = time.Minute * 30

func main() {
	users = append(users,
		User{Username: "admin", Password: "password", Role: "admin"},
		User{Username: "voter1", Password: "voter1", Role: "voter"},
		User{Username: "voter2", Password: "voter2", Role: "voter"},
	)
	go expireSessions()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/poll", pollHandler)
	http.HandleFunc("/createPoll", createPollHandler)
	http.HandleFunc("/deletePoll", deletePollHandler)
	http.HandleFunc("/vote", voteHandler)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Same logic as before
	token := generateToken()
	user.Token = token
	user.LastActive = time.Now()
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(sessionExpirationDuration),
	})
	// Rest of the login logic
}

func getUserFromToken(r *http.Request) *User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	token := cookie.Value

	for _, user := range users {
		if user.Token == token {
			// Check if the session is expired
			if time.Since(user.LastActive) > sessionExpirationDuration {
				return nil // Session expired
			}
			// Update the last active time to reset the expiration timer
			user.LastActive = time.Now()
			return &user
		}
	}
	return nil
}

func expireSessions() {
	for {
		time.Sleep(time.Minute) // Check for expired sessions every minute
		for i, user := range users {
			if time.Since(user.LastActive) > sessionExpirationDuration {
				user.Token = "" // Invalidate the session token
				users[i] = user // Update the user struct in the slice
			}
		}
	}
}

// Rest of the handlers remain the same
