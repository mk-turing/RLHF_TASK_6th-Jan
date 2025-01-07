package main

import (
	"fmt"
)

// ErrorDSL defines the DSL for network errors.
type ErrorDSL struct {
	Protocol string // Protocol (e.g., "HTTP", "TCP")
	Type     string // Error type (e.g., "ConnectionError", "TimeoutError")
	Code     int    // Error code
	Message  string // Error message
	Context  string // Additional context
}

// Error returns an error object for the DSL.
func (e *ErrorDSL) Error() string {
	return fmt.Sprintf("%s: %s: %d - %s (Context: %s)", e.Protocol, e.Type, e.Code, e.Message, e.Context)
}

// HTTPConnectionError creates an HTTPConnectionError instance.
func HTTPConnectionError(message, context string) error {
	return &ErrorDSL{
		Protocol: "HTTP",
		Type:     "ConnectionError",
		Code:     201,
		Message:  message,
		Context:  context,
	}
}

// HTTPTimeoutError creates an HTTPTimeoutError instance.
func HTTPTimeoutError(message, context string) error {
	return &ErrorDSL{
		Protocol: "HTTP",
		Type:     "TimeoutError",
		Code:     202,
		Message:  message,
		Context:  context,
	}
}

// TCPConnectionError creates a TCPConnectionError instance.
func TCPConnectionError(message, context string) error {
	return &ErrorDSL{
		Protocol: "TCP",
		Type:     "ConnectionError",
		Code:     101,
		Message:  message,
		Context:  context,
	}
}

// TCPTimeoutError creates a TCPTimeoutError instance.
func TCPTimeoutError(message, context string) error {
	return &ErrorDSL{
		Protocol: "TCP",
		Type:     "TimeoutError",
		Code:     102,
		Message:  message,
		Context:  context,
	}
}

// HTTPConnect establishes an HTTP connection to a remote server.
func HTTPConnect(host, port string) error {
	// Simulated HTTP network error
	return HTTPTimeoutError("Failed to connect to HTTP server", fmt.Sprintf("%s:%s", host, port))
}

// TCPConnect establishes a TCP connection to a remote server.
func TCPConnect(host, port string) error {
	// Simulated TCP network error
	return TCPConnectionError("Failed to connect to TCP server", fmt.Sprintf("%s:%s", host, port))
}

func main() {
	host := "example.com"
	port := "8080"

	// Attempt to connect using HTTP
	err := HTTPConnect(host, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("HTTP connection established successfully.")

	// Attempt to connect using TCP
	err = TCPConnect(host, port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("TCP connection established successfully.")
}
