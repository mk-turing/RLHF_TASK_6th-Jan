package main

import (
	"fmt"
	"sync"
	"time"
)

// ErrorLog is a map to store logged errors for analysis.
var ErrorLog = make(map[string][]ErrorLogEntry, 0)
var ErrorLogLock sync.Mutex

type ErrorDSL struct {
	Protocol string // Protocol (e.g., "HTTP", "TCP")
	Type     string // Error type (e.g., "ConnectionError", "TimeoutError")
	Code     int    // Error code
	Message  string // Error message
	Context  string // Additional context
}

func (e *ErrorDSL) Error() string {
	return fmt.Sprintf("%s: %s: %d - %s (Context: %s)", e.Protocol, e.Type, e.Code, e.Message, e.Context)
}

// ErrorLogEntry represents a single logged error entry.
type ErrorLogEntry struct {
	Error     ErrorDSL
	Timestamp time.Time
}

// LogError logs an error using the DSL.
func LogError(err error) {
	if dslErr, ok := err.(ErrorDSL); ok {
		ErrorLogLock.Lock()
		defer ErrorLogLock.Unlock()
		ErrorLog[dslErr.Protocol+"-"+dslErr.Type] = append(ErrorLog[dslErr.Protocol+"-"+dslErr.Type], ErrorLogEntry{
			Error:     dslErr,
			Timestamp: time.Now(),
		})
		fmt.Printf("Logged Error: %s\n", dslErr)
	} else {
		fmt.Printf("Unrecognized error type: %v\n", err)
	}
}

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
		LogError(err)
		return
	}
	fmt.Println("HTTP connection established successfully.")

	// Attempt to connect using TCP
	err = TCPConnect(host, port)
	if err != nil {
		LogError(err)
		return
	}
	fmt.Println("TCP connection established successfully.")

	// Example of analyzing errors
	AnalyzeErrors()
}

// AnalyzeErrors generates actionable insights from the logged errors.
func AnalyzeErrors() {
	ErrorLogLock.Lock()
	defer ErrorLogLock.Unlock()

	for key, entries := range ErrorLog {
		fmt.Printf("Protocol: %s, Type: %s\n", key[:key.Index("-")], key[key.Index("-")+1:])
		fmt.Printf("Total Errors: %d\n", len(entries))
		for _, entry := range entries {
			fmt.Printf("Time: %s, Error: %s\n", entry.Timestamp.Format(time.RFC3339), entry.Error)
		}
		fmt.Println()
	}
}
