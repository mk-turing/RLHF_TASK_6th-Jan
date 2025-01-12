package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
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
func (l *Library) GetBookByISBN(isbn string) (*Book, error) {
	for _, book := range l.Books {
		if book.ISBN == isbn {
			return &book, nil
		}
	}
	return nil, errors.New("book not found")
}

// BorrowBook borrows a book by its ISBN.
func (l *Library) BorrowBook(isbn string) error {
	book, err := l.GetBookByISBN(isbn)
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
	book, err := l.GetBookByISBN(isbn)
	if err != nil {
		return err
	}
	if book.Available {
		return errors.New("book is already available")
	}
	book.Available = true
	return nil
}

// RemoveBook removes a book by its ISBN.
func (l *Library) RemoveBook(isbn string) error {
	for i, book := range l.Books {
		if book.ISBN == isbn {
			l.Books = append(l.Books[:i], l.Books[i+1:]...)
			return nil
		}
	}
	return errors.New("book not found")
}

// UpdateBook updates a book's details by its ISBN.
func (l *Library) UpdateBook(isbn string, title, author, newISBN string) error {
	book, err := l.GetBookByISBN(isbn)
	if err != nil {
		return err
	}
	if title != "" {
		book.Title = title
	}
	if author != "" {
		book.Author = author
	}
	if newISBN != "" {
		book.ISBN = newISBN
		for _, b := range l.Books {
			if b.ISBN == newISBN && b.ISBN != isbn {
				return errors.New("duplicate ISBN found")
			}
		}
	}
	return nil
}

// SearchBooks searches for books based on the given attribute and value.
func (l *Library) SearchBooks(attribute, value string) []Book {
	var results []Book
	for _, book := range l.Books {
		switch attribute {
		case "title":
			if strings.Contains(strings.ToLower(book.Title), strings.ToLower(value)) {
				results = append(results, book)
			}
		case "author":
			if strings.Contains(strings.ToLower(book.Author), strings.ToLower(value)) {
				results = append(results, book)
			}
		case "isbn":
			if strings.Contains(book.ISBN, value) {
				results = append(results, book)
			}
		default:
			fmt.Println("Invalid search attribute. Please try again.")
			return nil
		}
	}
	return results
}

// SortBooks sorts the books based on the given attribute.
func (l *Library) SortBooks(attribute string) {