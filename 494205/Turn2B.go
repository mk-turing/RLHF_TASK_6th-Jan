package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type ErrorDSL struct {
	Type    string // Error type (e.g., "HTTPError", "TCPError", "GenericError")
	Code    int    // Error code (e.g., HTTP status code for HTTP, TCP error code for TCP)
	Message string // Error message
	Protocol string // Protocol associated with the error ("HTTP", "TCP", "")
	Context string // Additional context
}

func (e *ErrorDSL) Error() string {
	return fmt.Sprintf("%s: %d - %s (Protocol: %s, Context: %s)", e.Type, e.Code, e.Message, e.Protocol, e.Context)
}

// HTTP-specific error functions
func HTTPError(statusCode int, message, context string) error {
	return &ErrorDSL{
		Type:    "HTTPError",
		Code:    statusCode,
		Message: message,
		Protocol: "HTTP",
		Context: context,
	}
}

func TimeoutHTTPError(method, url string, duration time.Duration) error {
	return HTTPError(http.StatusRequestTimeout, "HTTP request timed out", fmt.Sprintf("%s to %s after %s", method, url, duration.String()))
}

// TCP-specific error functions
func TCPError(code int, message, context string) error {
	return &ErrorDSL{
		Type:    "TCPError",
		Code:    code,
		Message: message,
		Protocol: "TCP",
		Context: context,
	}
}

func TimeoutTCPError(host, port string, duration time.Duration) error {
	return TCPError(int(net.ERR_CONNECTION_TIMEOUT.Temporary()), "TCP connection timed out", fmt.Sprintf("Connection to %s:%s after %s", host, port, duration.String()))
}

// Generic error function
func GenericError(message, context string) error {
	return &ErrorDSL{
		Type:    "GenericError",
		Code:    0, // No specific code for generic errors
		Message: message,
		Context: context,
	}
}

// MakeHTTPClient creates an HTTP client with custom timeouts for testing.
func MakeHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

// SendHTTPRequest sends an HTTP request to the specified URL.
func SendHTTPRequest(client *http.Client, method, url string) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return GenericError("Failed to create HTTP request", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return TimeoutHTTPError(method, url, client.Timeout)
		}
		return GenericError("Failed to send HTTP request", err.Error())
	}
	defer resp.Body.Close()

	// Simulate an HTTP error for testing
	return HTTPError(500, "Internal Server Error", "Custom error from SendHTTPRequest")
}

// EstablishTCPConnection establishes a TCP connection to the specified host and port.
func EstablishTCPConnection(host, port string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		if errors.Is(err, net.ErrDialTimeout) {
			return TimeoutTCPError(host, port, timeout)
		}
		return TCPError(int(err.(net.Error).Temporary()), "Failed to establish TCP connection", err.Error())
	}
	defer conn.Close()