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
func EventHandler(dataStore *sync.Map, updateChan chan Event, processedEvents map[string]bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for event := range updateChan {
		// Check if the event has already been processed
		if processedEvents[event.Key] {
			continue
		}

		// Use LoadOrStore method to handle collisions automatically.
		oldValue, _ := dataStore.LoadOrStore(event.Key, event.Value)

		// Handle collisions by incrementing the value (assuming numeric values).
		if value, ok := oldValue.(int); ok {
			dataStore.Store(event.Key, value+event.Value)
		}

		// Mark the event as processed
		processedEvents[event.Key] = true
	}
}

func main() {
	// Create a synchronized map.
	dataStore := new(sync.Map)

	// Create a buffered channel to handle event updates.
	updateChan := make(chan Event, 100)

	// Create a map to store processed events
	processedEvents := make(map[string]bool)

	// Start multiple goroutines to handle event updates.
	const numHandlers = 5
	var wg sync.WaitGroup
	wg.Add(numHandlers)
	for i := 0; i < numHandlers; i++ {
		go EventHandler(dataStore, updateChan, processedEvents, &wg)
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
	wg.Wait()

	// Print the data store.
	dataStore.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: %d\n", key.(string), value.(int))
		return true
	})
}
