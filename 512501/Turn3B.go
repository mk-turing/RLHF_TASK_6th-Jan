package main

import (
	"fmt"
	"sync"
)

// Event represents an event that contains a key and a value.
type Event struct {
	Key   string
	Value int
}

// EventHandler processes events and updates the data store.
func EventHandler(dataStore *sync.Map, updateChan chan Event) {
	for event := range updateChan {
		// Use LoadOrStore method to handle collisions automatically.
		oldValue, _ := dataStore.LoadOrStore(event.Key, event.Value)

		// Handle collisions by incrementing the value (assuming numeric values).
		if value, ok := oldValue.(int); ok {
			dataStore.Store(event.Key, value+event.Value)
		}
	}
}

func main() {
	// Create a synchronized map.
	dataStore := new(sync.Map)

	// Create a buffered channel to handle event updates.
	updateChan := make(chan Event, 100)

	// Start multiple goroutines to handle event updates.
	const numHandlers = 5
	for i := 0; i < numHandlers; i++ {
		go EventHandler(dataStore, updateChan)
	}

	// Generate and send events.
	const numEvents = 10000
	for i := 0; i < numEvents; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := i + 1
		updateChan <- Event{Key: key, Value: value}

		// Simulate simultaneous events by triggering the same key multiple times.
		for j := 0; j < 3; j++ {
			updateChan <- Event{Key: key, Value: value}
		}
	}

	// Close the update channel to signal the handlers to finish processing.
	close(updateChan)

	// Wait for all event handlers to finish.
	for i := 0; i < numHandlers; i++ {
		<-updateChan
	}

	// Print the data store.
	dataStore.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: %d\n", key.(string), value.(int))
		return true
	})
}
