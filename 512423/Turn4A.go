package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Define a schema type
type Schema struct {
	Version int
	Fields  []Field
}

type Field struct {
	Name    string
	Type    string
	Default interface{} `json:",omitempty"`
}

// Simulated schema registry
type SchemaRegistry struct {
	schemas map[int]Schema
	mu      sync.RWMutex
}

func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		schemas: make(map[int]Schema),
	}
}

func (registry *SchemaRegistry) RegisterSchema(version int, schema Schema) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.schemas[version] = schema
}

func (registry *SchemaRegistry) GetSchema(version int) (Schema, error) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	schema, ok := registry.schemas[version]
	if !ok {
		return Schema{}, fmt.Errorf("schema not found for version %d", version)
	}
	return schema, nil
}

// Simulated message queue message
type WorkMessage struct {
	BatchID       int
	Batch         [][]byte
	SchemaVersion int
}

type ResultMessage struct {
	BatchID       int
	Results       []byte
	SchemaVersion int
}

// Simulated message queue interface
type MessageQueue interface {
	Send(ctx context.Context, msg interface{}) error
	Receive(ctx context.Context) (interface{}, error)
}

func extractData(batchSize int, registry *SchemaRegistry) [][]byte {
	data := [][]byte{
		[]byte(`{"id": 1, "name": "Alice", "age": 25, "schemaVersion": 1}`),
		[]byte(`{"id": 2, "name": "Bob", "age": 30, "schemaVersion": 2}`),
		[]byte(`{"id": 3, "name": "Charlie", "age": 18, "schemaVersion": 2}`),
	}

	batches := make([][]byte, 0)
	for i := 0; i < len(data); i += batchSize {
		end := min(i+batchSize, len(data))
		// Append each item in the slice, not the whole slice
		batches = append(batches, data[i:end]...)
	}
	return batches
}

// Function to transform data from the batch
func transformData(workMessage WorkMessage, registry *SchemaRegistry) []byte {
	var transformedRecords []map[string]interface{}
	schema, err := registry.GetSchema(workMessage.SchemaVersion)
	if err != nil {
		fmt.Println("Error fetching schema:", err)
		return nil
	}

	// Process each record in the batch
	for _, recordBytes := range workMessage.Batch {
		var record map[string]interface{}
		err := json.Unmarshal(recordBytes, &record)
		if err != nil {
			fmt.Println("Error unmarshaling record:", err)
			continue
		}

		// Create transformed record based on the schema
		transformedRecord := make(map[string]interface{})
		for _, field := range schema.Fields {
			value, ok := record[field.Name]
			if !ok {
				value = field.Default
			}
			transformedRecord[field.Name] = value
		}
		transformedRecords = append(transformedRecords, transformedRecord)
	}

	// Marshal the array of transformed records into a JSON array
	resultBytes, err := json.Marshal(transformedRecords)
	if err != nil {
		fmt.Println("Error marshaling transformed records:", err)
		return nil
	}

	// Debugging: Print the raw JSON output
	fmt.Printf("Transformed Data (JSON): %s\n", string(resultBytes))

	return resultBytes
}

// Function to load and print the transformed data
func loadData(results []byte, schemaVersion int) {
	// Debugging: Print the raw JSON data
	fmt.Printf("Raw JSON Data: %s\n", string(results))

	var transformedData []map[string]interface{}
	err := json.Unmarshal(results, &transformedData)
	if err != nil {
		fmt.Println("Error loading data:", err)
		return
	}

	// Print each item from the unmarshaled data
	for _, data := range transformedData {
		fmt.Printf("Loading data: %v\n", data)
	}
}

func main() {
	batchSize := 1
	schemaRegistry := NewSchemaRegistry()
	schemaRegistry.RegisterSchema(1, Schema{
		Version: 1,
		Fields: []Field{
			{"id", "int", nil},
			{"name", "string", nil},
			{"age", "int", nil},
			{"schemaVersion", "int", nil},
		},
	})
	schemaRegistry.RegisterSchema(2, Schema{
		Version: 2,
		Fields: []Field{
			{"id", "int", nil},
			{"name", "string", nil},
			{"age", "int", nil},
			{"schemaVersion", "int", nil},
			{"email", "string", ""},
		},
	})

	batches := extractData(batchSize, schemaRegistry)
	var wg sync.WaitGroup

	ctx := context.Background()

	mq := NewSimulatedMessageQueue()

	// Start message queue transfer goroutine
	go mq.StartTransfer()

	// Distribute work
	for i, batch := range batches {
		wg.Add(1)
		go func(batchID int, batch [][]byte) {
			defer wg.Done()

			workMessage := WorkMessage{
				BatchID:       batchID,
				Batch:         batch,
				SchemaVersion: 1,
			}
			if err := mq.Send(ctx, workMessage); err != nil {
				fmt.Println("Error sending work message:", err)
				return
			}
		}(i, [][]byte{batch})
	}

	// Wait for all messages to be sent
	wg.Wait()
	close(mq.sendChan) // Signal no more messages will be sent

	// Collect results
	var allResults []byte
	for {
		msg, err := mq.Receive(ctx)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			break
		}
		if msg == nil {
			break
		}

		resultMessage := msg.(WorkMessage)
		results := transformData(resultMessage, schemaRegistry)
		allResults = append(allResults, results...)
	}

	// Load all transformed data
	loadData(allResults, 1)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type SimulatedMessageQueue struct {
	sendChan    chan interface{}
	receiveChan chan interface{}
}

func NewSimulatedMessageQueue() *SimulatedMessageQueue {
	return &SimulatedMessageQueue{
		sendChan:    make(chan interface{}, 10), // Buffered channels to avoid blocking
		receiveChan: make(chan interface{}, 10),
	}
}

func (mq *SimulatedMessageQueue) StartTransfer() {
	for msg := range mq.sendChan {
		mq.receiveChan <- msg
	}
	close(mq.receiveChan) // Signal no more messages will be transferred
}

func (mq *SimulatedMessageQueue) Send(ctx context.Context, msg interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case mq.sendChan <- msg:
		return nil
	}
}

func (mq *SimulatedMessageQueue) Receive(ctx context.Context) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg, ok := <-mq.receiveChan:
		if !ok {
			return nil, nil // Channel closed, no more messages
		}
		return msg, nil
	}
}
