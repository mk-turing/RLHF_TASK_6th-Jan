package main

import (
	"fmt"
	"log"
)

type Book struct {
	Title  string
	Author string
	ID     string
}

func (b *Book) AddBook(newBook Book) {
	// Append a new book to the collection
	b.Book = append(b.Book, newBook)
}

type User struct {
	Name string
	ID   string
}

func (u *User) AddUser(newUser User) {
	// Append a new user to the collection
	u.User = append(u.User, newUser)
}

func main() {
	// Initialize book and user slices
	books := []Book{}
	users := []User{}

	// Add some books and users for demonstration
	books = append(books, Book{Title: "Go Programming", Author: "Alan A. A. Donovan", ID: "123"})
	users = append(users, User{Name: "Alice", ID: "U101"})

	// Display the collections
	fmt.Println("Books:")
	for _, b := range books {
		fmt.Println(b)
	}

	fmt.Println("\nUsers:")
	for _, u := range users {
		fmt.Println(u)
	}

	log.Println("Library management system running.")
}
