package main

import (
	"fmt"
	"log"
	"time"
)

// ErrorDSL defines the DSL for network errors with logging fields.
type ErrorDSL struct {
	Protocol string    // Protocol (e.g., "HTTP", "TCP")
	Type     string    // Error type (e.g., "ConnectionError", "TimeoutError")
	Code     int       // Error code
	Message  string    // Error message
	Context  string    // Additional context
	Time     time.Time // Timestamp of the error
	Source   string    // Source of the error (e.g., function name)
}

// Error returns an error object for the DSL.
func (e *ErrorDSL) Error() string {
	return fmt.Sprintf("%s: %s: %d - %s (Context: %s)", e.Protocol, e.Type, e.Code, e.Message, e.Context)
}

// NewErrorDSL creates a new ErrorDSL instance with logging fields.
func NewErrorDSL(protocol, errorType string, code int, message, context, source string) error {
	return &ErrorDSL{
		Protocol: protocol,
		Type:     errorType,
		Code:     code,
		Message:  message,
		Context:  context,
		Time:     time.Now(),
		Source:   source,
	}
}

// customLogger defines a custom logger for logging errors using the DSL.
type customLogger struct {
	*log.Logger
}

// LogError logs an error using the DSL.
func (l *customLogger) LogError(err error) {
	if dslError, ok := err.(*ErrorDSL); ok {
		l.Printf("Error: %+v", dslError)
	} else {
		l.Printf("Error: %v", err)
	}
}

// HTTPConnect establishes an HTTP connection to a remote server.
func HTTPConnect(host, port string) error {
	// Simulated HTTP network error
	return NewErrorDSL("HTTP", "ConnectionError", 201, "Failed to connect to HTTP server", fmt.Sprintf("%s:%s", host, port), "HTTPConnect")
}

// TCPConnect establishes a TCP connection to a remote server.
func TCPConnect(host, port string) error {
	// Simulated TCP network error
	return NewErrorDSL("TCP", "TimeoutError", 102, "Failed to connect to TCP server", fmt.Sprintf("%s:%s", host, port), "TCPConnect")
}

func main() {
	logger := customLogger{log.New(log.Writer(), "", log.LstdFlags)}

	host := "example.com"
	port := "8080"

	// Attempt to connect using HTTP
	err := HTTPConnect(host, port)
	if err != nil {
		logger.LogError(err)
	}
	fmt.Println("HTTP connection established successfully.")

	// Attempt to connect using TCP
	err = TCPConnect(host, port)
	if err != nil {
		logger.LogError(err)
	}
	fmt.Println("TCP connection established successfully.")
}
