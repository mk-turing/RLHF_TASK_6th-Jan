package main

import (
	"errors"
	"fmt"
	"log"
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

func main() {
	// Create a new library
	lib := NewLibrary()

	// Add some books to the library
	lib.AddBook(Book{Title: "The Go Programming Language", Author: "Alan A. A. Donovan", ISBN: "9780134190440", Available: true})
	lib.AddBook(Book{Title: "Introduction to Algorithms", Author: "Thomas H. Cormen", ISBN: "9780262033848", Available: true})
	lib.AddBook(Book{Title: "Clean Code", Author: "Robert C. Martin", ISBN: "9780132350884", Available: true})

	// List available books
	fmt.Println("Available Books:")
	availableBooks := lib.ListAvailableBooks()
	for _, book := range availableBooks {
		fmt.Printf("- %s by %s (ISBN: %s)\n", book.Title, book.Author, book.ISBN)
	}

	// Borrow a book
	err := lib.BorrowBook("9780134190440")
	if err != nil {
		log.Println("Error borrowing book:", err)
	} else {
		fmt.Println("\nBorrowed 'The Go Programming Language'")
	}

	// Try to borrow the same book again
	err = lib.BorrowBook("9780134190440")
	if err != nil {
		log.Println("Error borrowing book:", err)
	}

	// Return the borrowed book
	err = lib.ReturnBook("9780134190440")
	if err != nil {
		log.Println("Error returning book:", err)
	} else {
		fmt.Println("\nReturned 'The Go Programming Language'")
	}

	// List available books again after borrowing and returning a book
	fmt.Println("\nAvailable Books after borrowing and returning a book:")
	availableBooks = lib.ListAvailableBooks()
	for _, book := range availableBooks {
		fmt.Printf("- %s by %s (ISBN: %s)\n", book.Title, book.Author, book.ISBN)
	}
}
