package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	Username string
	Password string
	Token    string
	Role     string
	Votes    []int // added a slice to track user votes
}

var users []User
var sessions []string
var polls []string // add polls slice to store poll questions

func main() {
	polls = append(polls, "Poll Question 1", "Poll Question 2", "Poll Question 3")
	// Initialize votes slice for each user with an empty array
	for i := range users {
		users[i].Votes = []int{}
	}

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

func pollHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromToken(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Welcome to the poll page, %s! <br>", user.Username)

	// Display all polls and allow user to vote
	for pollIndex, pollQuestion := range polls {
		fmt.Fprintf(w, "<form action='/vote' method='post'>")
		fmt.Fprintf(w, "<input type='hidden' name='pollIndex' value='%d'>", pollIndex)
		fmt.Fprintf(w, "<p>%s</p>", pollQuestion)
		fmt.Fprintf(w, "<input type='radio' name='option' value='1'> Option 1<br>")
		fmt.Fprintf(w, "<input type='radio' name='option' value='2'> Option 2<br>")
		fmt.Fprintf(w, "<input type='submit' value='Vote'>")
		fmt.Fprintf(w, "</form>")
	}
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := getUserFromToken(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	pollIndex, err := strconv.Atoi(r.FormValue("pollIndex"))
	if err != nil {
		http.Error(w, "Invalid poll index", http.StatusBadRequest)
		return
	}

	// Check if the user has already voted for this poll
	if contains(user.Votes, pollIndex) {
		http.Error(w, "You have already voted for this poll.", http.StatusForbidden)
		return
	}

	option := r.FormValue("option")
	if option != "1" && option != "2" {
		http.Error(w, "Invalid vote option", http.StatusBadRequest)
		return
	}

	// Track the vote by appending the poll index to the user's votes slice
	user.Votes = append(user.Votes, pollIndex)
	fmt.Fprintf(w, "Thank you for your vote!")
}

func contains(slice []int, target int) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Same logic as before
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
