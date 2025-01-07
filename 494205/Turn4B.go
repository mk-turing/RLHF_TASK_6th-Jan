package main
import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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

// errorHandlingRules stores the error handling rules.
var errorHandlingRules map[string]func(error) error

// reloadConfig reloads the configuration file and updates the error handling rules.
func reloadConfig() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
		return
	}

	newErrorHandlingRules := make(map[string]func(error) error)

	if err := viper.UnmarshalKey("error_handling_rules", &newErrorHandlingRules); err != nil {
		log.Printf("Error unmarshaling error handling rules: %s", err)
		return
	}

	errorHandlingRules = newErrorHandlingRules
}

// initializeConfig initializes the configuration and sets up the watcher.
func initializeConfig() {
	reloadConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		reloadConfig()
	})
}

// handleError handles an error using the appropriate error handling rule.
func handleError(err error) {
	if err == nil {
		return
	}

	dslError, ok := err.(*ErrorDSL)
	if !ok {
		log.Printf("Error: %v", err)
		return
	}

	if rule, ok := errorHandlingRules[dslError.Type]; ok {
		rule(err)
	} else {
		log.Printf("Error: %+v", dslError)
	}
}

// HTTPConnect establishes an HTTP connection to a remote server.
func HTTPConnect(host, port string) error {
	// Simulated HTTP network error
	return &ErrorDSL{
		Protocol: "HTTP",
		Type:    "ConnectionError",
		Code:    201,
		Message:  "Failed to connect to HTTP server",
		Context:  fmt.Sprintf("%s:%s", host, port),
	}
}

// TCPConnect establishes a TCP connection to a remote server.
func TCPConnect(host, port string) error {
	// Simulated TCP network error
	return &ErrorDSL{
		Protocol: "TCP",
		Type:    "TimeoutError",
		Code:    102,
		Message:  "Failed to connect to TCP server",
		Context:  fmt.Sprintf("%s:%s", host, port),
	}
}

func main() {
	initializeConfig()

	host := "example.com"
	port := "8080"

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()