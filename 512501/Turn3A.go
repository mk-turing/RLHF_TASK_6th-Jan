package main

import (
	"fmt"
	"sync"
	"time"
)

// Event represents an event that contains a key and a value
type Event struct {
	Key   string
	Value int
}

// EventHandler processes events and updates the data store
func EventHandler(dataStore map[string]*eventRecord, rwmu *sync.RWMutex, events chan Event) {
	for event := range events {
		// Lock for writing to the map
		rwmu.Lock()
		defer rwmu.Unlock()

		// Handle the event
		key := event.Key
		value := event.Value

		if record, exists := dataStore[key]; exists {
			// Handle collision: increment the count or overwrite the value
			record.Count += 1
			record.Total += value
		} else {
			// No collision: create a new record
			dataStore[key] = &eventRecord{Total: value, Count: 1}
		}
	}
}

// eventRecord holds the total value and count for a key
type eventRecord struct {
	Total int
	Count int
}

func main() {
	// Create a data store map
	dataStore := make(map[string]*eventRecord)

	// Create a channel to receive events
	events := make(chan Event, 1000) // Buffered channel to handle event bursts

	// Start multiple goroutines to handle events
	const numHandlers = 5
	for i := 0; i < numHandlers; i++ {
		go EventHandler(dataStore, &sync.RWMutex{}, events)
	}

	// Generate and send events
	const numEvents = 10000
	for i := 0; i < numEvents; i++ {
		key := fmt.Sprintf("key_%d", i)
		value := i + 1
		events <- Event{Key: key, Value: value}

		// Simulate simultaneous events by triggering the same key multiple times
		for j := 0; j < 3; j++ {
			events <- Event{Key: key, Value: value}
		}
	}

	// Graceful shutdown: wait for event handlers to consume all events
	time.Sleep(2 * time.Second) // Wait some time for events to be processed
	close(events)

	// Wait for all event handlers to finish
	for i := 0; i < numHandlers; i++ {
		go func() {
			<-events // Block until the channel is closed
		}()
	}

	// Print the data store
	var wg sync.WaitGroup
	for key := range dataStore {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			rwmu := &sync.RWMutex{}
			rwmu.RLock()
			record := dataStore[k]
			rwmu.RUnlock()
			avg := record.Total / record.Count
			fmt.Printf("%s: Total = %d, Count = %d, Average = %d\n", k, record.Total, record.Count, avg)
		}(key)
	}
	wg.Wait()
}
