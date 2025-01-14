package main

import (
	"fmt"
	"sync"
)

// Define an event type
type Event struct {
	Key   string
	Value string
}

// define a handler function type
type EventHandler func(*Event)

// we'll use a map to store events with a key-handler map as value.
type EventManager struct {
	Events map[string]map[EventHandler]bool
	mtx    *sync.RWMutex
}

// Create a new event manager
func NewEventManager() *EventManager {
	return &EventManager{
		Events: make(map[string]map[EventHandler]bool),
		mtx:    &sync.RWMutex{},
	}
}

// AddEventHandler will add an event handler to the event manager.
func (e *EventManager) AddEventHandler(key string, handler EventHandler) {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	if e.Events[key] == nil {
		e.Events[key] = make(map[EventHandler]bool)
	}
	e.Events[key][handler] = true
}

// TriggerEvent will trigger an event with given key and value.
func (e *EventManager) TriggerEvent(key, value string) {
	e.mtx.RLock()
	defer e.mtx.RUnlock()
	if handlers, ok := e.Events[key]; ok {
		event := &Event{Key: key, Value: value}
		for handler := range handlers {
			handler(event)
		}
	}
}

func main() {
	//create an event manager.
	manager := NewEventManager()

	// Create event handlers
	handler1 := func(event *Event) {
		fmt.Println("Handler 1 received event:", event.Key, ":", event.Value)
	}
	handler2 := func(event *Event) {
		fmt.Println("Handler 2 received event:", event.Key, ":", event.Value)
	}
	// Define 2 events
	eventKey1 := "event1"
	eventKey2 := "event2"
	// registering handler for event 1
	manager.AddEventHandler(eventKey1, handler1)
	manager.AddEventHandler(eventKey1, handler2)
	// registering handler for event 2
	manager.AddEventHandler(eventKey2, handler1)
	// Trigger events
	manager.TriggerEvent(eventKey1, "Event1 Data")
	manager.TriggerEvent(eventKey2, "Event2 Data")
	// Output will be
	/*Handler 1 received event: event1 : Event1 Data
	  Handler 2 received event: event1 : Event1 Data
	  Handler 1 received event: event2 : Event2 Data*/

}
