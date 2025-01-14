package main

import (
	"fmt"
	"sync"
)

// Event represents an event that contains a key and a value
type Event struct {
	Key   string
	Value int
}

// EventHandler processes events and updates the data store
func EventHandler(dataStore map[string]int, mu *sync.Mutex, wg *sync.WaitGroup, events chan Event) {
	defer wg.Done()

	for event := range events {
		// Handle the event
		key := event.Key
		value := event.Value

		// Lock the map to avoid concurrent writes
		mu.Lock()
		// Detect and manage key collisions
		collisionCount := 0
		for {
			_, exists := dataStore[key]
			if !exists {
				// Key does not exist, no collision
				dataStore[key] = value
				break
			}
			// Key exists, handle collision
			collisionCount++
			key = fmt.Sprintf("%s_%d", event.Key, collisionCount)
		}
		// Unlock after operation
		mu.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to protect map access

	// Create a data store map
	dataStore := make(map[string]int)

	// Create a channel to receive events
	events := make(chan Event)

	// Start multiple goroutines to handle events
	const numHandlers = 3
	for i := 0; i < numHandlers; i++ {
		wg.Add(1)
		go EventHandler(dataStore, &mu, &wg, events)
	}

	// Generate and send events
	const numEvents = 10
	for i := 0; i < numEvents; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := i + 1
		events <- Event{Key: key, Value: value}

		// Simulate simultaneous events by triggering the same key multiple times
		for j := 0; j < 3; j++ {
			events <- Event{Key: key, Value: value}
		}
	}

	// Close the event channel
	close(events)

	// Wait for all event handlers to finish
	wg.Wait()

	// Print the data store
	fmt.Println("Data Store:")
	for key, value := range dataStore {
		fmt.Printf("%s: %d\n", key, value)
	}
}
