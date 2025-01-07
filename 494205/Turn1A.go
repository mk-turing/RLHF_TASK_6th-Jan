package main

import (
	"fmt"
)

// ErrorDSL defines the DSL for network errors.
type ErrorDSL struct {
	Type    string // Error type (e.g., "ConnectionError", "TimeoutError")
	Code    int    // Error code
	Message string // Error message
	Context string // Additional context
}

// Error returns an error object for the DSL.
func (e *ErrorDSL) Error() string {
	return fmt.Sprintf("%s: %d - %s (Context: %s)", e.Type, e.Code, e.Message, e.Context)
}

// ConnectionError creates a ConnectionError instance.
func ConnectionError(message, context string) error {
	return &ErrorDSL{
		Type:    "ConnectionError",
		Code:    101,
		Message: message,
		Context: context,
	}
}

// TimeoutError creates a TimeoutError instance.
func TimeoutError(message, context string) error {
	return &ErrorDSL{
		Type:    "TimeoutError",
		Code:    102,
		Message: message,
		Context: context,
	}
}

// Connect establishes a connection to a remote server.
func Connect(host, port string) error {
	// Simulated network error
	return ConnectionError("Failed to connect to server", fmt.Sprintf("%s:%s", host, port))
}

func main() {
	host := "example.com"
	port := "8080"

	err := Connect(host, port)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connection established successfully.")
}
