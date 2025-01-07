package main

import "fmt"

// Define the ErrorCodes enum
type ErrorCodes int

const (
	ConnectionError ErrorCodes = iota
	TimeoutError
	NetworkError
	InvalidRequestError
	UnknownError
)

// Define the 'NetworkError' struct
type NetworkError struct {
	code ErrorCodes
	msg  string
}

func (ne *NetworkError) Error() string {
	return fmt.Sprintf("network error: %s", ne.msg)
}
func (ne *NetworkError) Code() ErrorCodes {
	return ne.code
}

// Create a 'NewNetworkError' function to return a new 'NetworkError' instance
func NewNetworkError(code ErrorCodes, msg string) error {
	return &NetworkError{code: code, msg: msg}
}

// Custom DSL functions
func IsConnectionError(err error) bool {
	if ne, ok := err.(*NetworkError); ok {
		return ne.Code() == ConnectionError
	}
	return false
}

func IsTimeoutError(err error) bool {
	if ne, ok := err.(*NetworkError); ok {
		return ne.Code() == TimeoutError
	}
	return false
}

// Your network application logic here
func NetworkOperation() error {
	// Simulate some network operation
	return NewNetworkError(TimeoutError, "timed out")
}

func main() {
	if err := NetworkOperation(); err != nil {
		if IsTimeoutError(err) {
			fmt.Println("Timeout occurred:", err)
		} else if IsConnectionError(err) {
			fmt.Println("Connection error:", err)
		} else {
			fmt.Println("Unknown error:", err)
		}
	}
}
