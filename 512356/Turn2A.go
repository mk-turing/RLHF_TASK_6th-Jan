package main

import (
	"fmt"
	"sync"
	"time"
)

// Define the logMessage type
type logMessage func(level, message string)

// Global variables to store log level info
var totalLogs int = 0
var levelLimit int = 5 // Limit logs to 5 for simplicity

// Mutex to synchronize access to shared variables
var mutex sync.Mutex

// A function that uses the global variable for counting and limiting logs
var logProcessor logMessage = func(level, message string) {
	mutex.Lock()
	defer mutex.Unlock()

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

// Define a custom processor function
func customLogProcessor(level, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	if message == "Error" {
		fmt.Printf("[ERROR]: %s\n", message)
	}
	// Modify the shared variable
	levelLimit++
}

func main() {
	// Register the custom processor
	RegisterLogProcessor(customLogProcessor)

	// Start multiple goroutines to log messages
	for i := 0; i < 10; i++ {
		go logProcessor("Info", fmt.Sprintf("This is message %d", i))
		go logProcessor("Error", "Something went wrong.")
	}

	// Wait for all goroutines to finish
	time.Sleep(2 * time.Second)
}
