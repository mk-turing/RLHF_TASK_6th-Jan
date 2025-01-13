package main

import (
	"fmt"
	"time"
)

// Define the logMessage type
type logMessage func(level, message string)

// Global variable to store log level info
var totalLogs int = 0
var levelLimit int = 5 // Limit logs to 5 for simplicity

// A function that uses the global variable for counting and limiting logs
var logProcessor logMessage = func(level, message string) {
	totalLogs++
	if totalLogs <= levelLimit {
		fmt.Printf("%v: %v\n", time.Now(), message)
	}
}

// Register a log processing function
func RegisterLogProcessor(processor logMessage) {
	// Assign the processor function to logProcessor
	logProcessor = processor
}

func main() {
	// Use the logger
	logProcessor("Info", "This is an informational message.")
	logProcessor("Info", "Another informational message.")

	// Define a custom processor function
	customLogProcessor := func(level, message string) {
		if message == "Error" {
			fmt.Printf("[ERROR]: %s\n", message)
		}
		// Modify the shared variable
		levelLimit++
	}

	// Register the custom processor
	RegisterLogProcessor(customLogProcessor)

	// Log another message
	logProcessor("Error", "Something went wrong.")
	logProcessor("Info", "Yet another informational message.")
}
