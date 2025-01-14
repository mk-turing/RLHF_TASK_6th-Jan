package main

import (
	"context"
	"fmt"
	"sync"
)

type SourceData struct {
	ID   int
	Name string
	Age  int
}

type TransformedData struct {
	ID       int
	Name     string
	AgeGroup string
	IsAdult  bool
}

// Simulated message queue message
type WorkMessage struct {
	BatchID int
	Batch   []SourceData
}

type ResultMessage struct {
	BatchID int
	Results []TransformedData
}

// Simulated message queue interface
type MessageQueue interface {
	Send(ctx context.Context, msg interface{}) error
	Receive(ctx context.Context) (interface{}, error)
}

func extractData(batchSize int) [][]SourceData {
	sourceData := []SourceData{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 18},
	}

	batches := make([][]SourceData, 0)
	for i := 0; i < len(sourceData); i += batchSize {
		batches = append(batches, sourceData[i:min(i+batchSize, len(sourceData))])
	}
	return batches
}

func transformData(workMessage WorkMessage) []TransformedData {
	var results []TransformedData
	for _, record := range workMessage.Batch {
		ageGroup := "Young"
		isAdult := false
		if record.Age >= 18 {
			ageGroup = "Adult"
			isAdult = true
		}

		results = append(results, TransformedData{
			ID:       record.ID,
			Name:     record.Name,
			AgeGroup: ageGroup,
			IsAdult:  isAdult,
		})
	}
	return results
}

func loadData(results []TransformedData) {
	for _, data := range results {
		fmt.Printf("Loading data: %v\n", data)
	}
}

func main() {
	batchSize := 2
	batches := extractData(batchSize)
	var wg sync.WaitGroup

	ctx := context.Background()

	mq := NewSimulatedMessageQueue()

	// Start message queue transfer goroutine
	go mq.StartTransfer()

	// Distribute work
	for i, batch := range batches {
		wg.Add(1)
		go func(batchID int, batch []SourceData) {
			defer wg.Done()

			workMessage := WorkMessage{
				BatchID: batchID,
				Batch:   batch,
			}
			if err := mq.Send(ctx, workMessage); err != nil {
				fmt.Println("Error sending work message:", err)
				return
			}
		}(i, batch)
	}

	// Wait for all messages to be sent
	wg.Wait()
	close(mq.sendChan) // Signal no more messages will be sent

	// Collect results
	var allResults []TransformedData
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
		results := transformData(resultMessage)
		allResults = append(allResults, results...)
	}

	// Load all transformed data
	loadData(allResults)
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
