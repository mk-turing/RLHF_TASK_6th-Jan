package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Define a User struct
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Define a Session struct to manage sessions
type Session struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
	Expiry int64  `json:"expiry"`
}

var users []User
var mu sync.RWMutex

func authenticateUser(username, password string) (bool, int, *User) {
	mu.RLock()
	defer mu.RUnlock()

	for _, user := range users {
		if user.Username == username && user.Password == password {
			return true, http.StatusOK, &user
		}
	}
	return false, http.StatusUnauthorized, nil
}

var sessions []Session
var sessionsMu sync.RWMutex

func createSession(userID int) string {
	token := generateToken() // Implement a function to generate a random token
	session := Session{
		UserID: userID,
		Token:  token,
		Expiry: getCurrentTime() + 3600, // Expires in one hour
	}

	sessionsMu.Lock()
	defer sessionsMu.Unlock()

	sessions = append(sessions, session)
	return token
}

func validateSession(token string) (bool, int) {
	sessionsMu.RLock()
	defer sessionsMu.RUnlock()

	for _, session := range sessions {
		if session.Token == token && session.Expiry > getCurrentTime() {
			return true, http.StatusOK
		}
	}
	return false, http.StatusUnauthorized
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if authenticated, status, user := authenticateUser(username, password); authenticated {
		token := createSession(user.ID)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprintln(w, map[string]string{"token": token})
	} else {
		w.WriteHeader(status)
		fmt.Fprintln(w, "Invalid credentials")
	}
}

func pollHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	if authenticated, status := validateSession(token); authenticated {
		w.WriteHeader(status)
		fmt.Fprintln(w, "Access granted to poll")
	} else {
		w.WriteHeader(status)
		fmt.Fprintln(w, "Unauthorized access")
	}
}

func main() {
	users = []User{
		{ID: 1, Username: "user1", Password: "hashedPassword1"},
		{ID: 2, Username: "user2", Password: "hashedPassword2"},
	}

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/poll", pollHandler)
	http.ListenAndServe(":8080", nil)
}
