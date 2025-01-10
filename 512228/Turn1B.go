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
}

var users []User
var sessions []string

func main() {
	users = append(users, User{Username: "admin", Password: "password"})
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/poll", pollHandler)
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
			// Generate a new token and store it in the user's struct and sessions slice
			token := generateToken()
			user.Token = token
			sessions = append(sessions, token)
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   token,
				Expires: time.Now().Add(30 * time.Minute), // Set token expiration time
			})
			http.Redirect(w, r, "/poll", http.StatusSeeOther)
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

	// Validate the token using the sessions slice
	for _, session := range sessions {
		if session == token {
			// Authenticated user, display poll page
			fmt.Fprintf(w, "Welcome to the poll page!")
			return
		}
	}
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func generateToken() string {
	// Generate a unique token using a cryptographically secure random number generator (CRNG)
	// For simplicity, we'll use a random string in this example
	return "random-token-string"
}
