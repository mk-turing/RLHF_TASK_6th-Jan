package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Book represents a book in the library.
type Book struct {
	Title     string
	Author    string
	ISBN      string
	Available bool
}

// Library represents a collection of books.
type Library struct {
	Books []Book
}

// NewLibrary creates a new empty library.
func NewLibrary() *Library {
	return &Library{Books: make([]Book, 0)}
}

// AddBook adds a new book to the library.
func (l *Library) AddBook(book Book) {
	l.Books = append(l.Books, book)
}

// GetBookByISBN retrieves a book by its ISBN.
func (l *Library) GetBookByISBN(isbn string) (*Book, int, error) {
	for index, book := range l.Books {
		if book.ISBN == isbn {
			return &book, index, nil
		}
	}
	return nil, -1, errors.New("book not found")
}

// RemoveBook removes a book from the library by its ISBN.
func (l *Library) RemoveBook(isbn string) error {
	book, index, err := l.GetBookByISBN(isbn)
	if err != nil {
		return err
	}
	l.Books = append(l.Books[:index], l.Books[index+1:]...)
	return nil
}

// UpdateBook updates a book in the library by its ISBN.
func (l *Library) UpdateBook(isbn string, updatedBook Book) error {
	book, index, err := l.GetBookByISBN(isbn)
	if err != nil {
		return err
	}
	l.Books[index] = updatedBook
	return nil
}

// BorrowBook borrows a book by its ISBN.
func (l *Library) BorrowBook(isbn string) error {
	book, _, err := l.GetBookByISBN(isbn)
	if err != nil {
		return err
	}
	if !book.Available {
		return errors.New("book is already borrowed")
	}
	book.Available = false
	return nil
}

// ReturnBook returns a book by its ISBN.
func (l *Library) ReturnBook(isbn string) error {
	book, _, err := l.GetBookByISBN(isbn)
	if err != nil {
		return err
	}
	if book.Available {
		return errors.New("book is already available")
	}
	book.Available = true
	return nil
}

// ListAvailableBooks lists all the available books in the library.
func (l *Library) ListAvailableBooks() []Book {
	var availableBooks []Book
	for _, book := range l.Books {
		if book.Available {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func getUserInput(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	scanner.Scan()
	return scanner.Text()
}

func validateIsbn(isbn string) error {
	// Simple validation to check for correct ISBN length
	if len(isbn) != 10 && len(isbn) != 13 {
		return errors.New("ISBN must be 10 or 13 characters long")
	}
	return nil
}

func main() {
	// Create a new library
	lib := NewLibrary()

	// User menu loop
	for {
		fmt.Println("\nLibrary Management System Menu:")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Update Book")
		fmt.Println("4. Borrow Book")
		fmt.Println("5. Return Book")
		fmt.Println("6. List Available Books")