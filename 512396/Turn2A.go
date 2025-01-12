package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
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
		// Check if the new ISBN is unique
		for _, b := range l.Books {
			if b.ISBN == newISBN && b.ISBN != isbn {
				return errors.New("duplicate ISBN found")
			}
		}
	}

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

// Main function to interact with the user.
func main() {
	lib := NewLibrary()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Library Management System!")

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Update Book")
		fmt.Println("4. Borrow Book")
		fmt.Println("5. Return Book")
		fmt.Println("6. List Available Books")
		fmt.Println("7. Exit")

		choice, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading choice:", err)
			break
		}

		switch choice {
		case "1\n":
			fmt.Print("Enter book title: ")
			title, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading title:", err)
				continue
			}

			fmt.Print("Enter book author: ")
			author, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading author:", err)
				continue
			}

			fmt.Print("Enter book ISBN: ")
			isbn, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading ISBN:", err)
				continue
			}

			available := true
			fmt.Print("Is the book available? (yes/no): ")
			availability, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading availability:", err)
				continue
			}

			if availability == "no\n" || availability == "No\n" {
				available = false
			}

			book := Book{Title: title, Author: author, ISBN: isbn, Available: available}
			lib.AddBook(book)
			fmt.Println("Book added successfully!")

		case "2\n":
			fmt.Print("Enter book ISBN to remove: ")
			isbn, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading ISBN:", err)
				continue
			}

			err = lib.RemoveBook(isbn)
			if err != nil {
				fmt.Println("Error removing book:", err)
			} else {
				fmt.Println("Book removed successfully!")
			}

		case "3\n":
			fmt.Print("Enter book ISBN to update: ")
			isbn, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading ISBN:", err)
				continue
			}

			var title, author, newISBN string
			fmt.Print("Enter new title (empty to keep current): ")
			title, err = reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading title:", err)
				continue
			}

			fmt.Print("Enter new author (empty to keep current): ")
			author, err = reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading author:", err)
				continue
			}

			fmt.Print("Enter new ISBN (empty to keep current): ")
			newISBN, err = reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading ISBN:", err)
				continue
			}

			err = lib.UpdateBook(isbn, title, author, newISBN)
			if err != nil {
				fmt.Println("Error updating book:", err)
			} else {
				fmt.Println("Book updated successfully!")
			}

		case "4\n":
			fmt.Print("Enter book ISBN to borrow: ")
			isbn, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading ISBN:", err)
				continue
			}

			err = lib.BorrowBook(isbn)
			if err != nil {
				fmt.Println("Error borrowing book:", err)
			} else {
				fmt.Println("Book borrowed successfully!")
			}

		case "5\n":
			fmt.Print("Enter book ISBN to return: ")
			isbn, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Error reading ISBN:", err)
				continue
			}

			err = lib.ReturnBook(isbn)
			if err != nil {
				fmt.Println("Error returning book:", err)
			} else {
				fmt.Println("Book returned successfully!")
			}

		case "6\n":
			fmt.Println("Available Books:")
			availableBooks := lib.ListAvailableBooks()
			for _, book := range availableBooks {
				fmt.Printf("- %s by %s (ISBN: %s)\n", book.Title, book.Author, book.ISBN)
			}

		case "7\n":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
