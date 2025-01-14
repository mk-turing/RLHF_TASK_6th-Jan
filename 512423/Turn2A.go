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

func extractData(batchSize int) [][]SourceData {
	// Simulate extracting data in batches
	sourceData := []SourceData{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 18},
		// Add more data as needed
	}

	batches := make([][]SourceData, 0)
	for i := 0; i < len(sourceData); i += batchSize {
		batches = append(batches, sourceData[i:min(i+batchSize, len(sourceData))])
	}
	return batches
}

func transformData(ctx context.Context, sourceData []SourceData, wg *sync.WaitGroup, results *[]TransformedData) {
	defer wg.Done()

	for _, record := range sourceData {
		ageGroup := "Young"
		isAdult := false
		if record.Age >= 18 {
			ageGroup = "Adult"
			isAdult = true
		}

		*results = append(*results, TransformedData{
			ID:       record.ID,
			Name:     record.Name,
			AgeGroup: ageGroup,
			IsAdult:  isAdult,
		})
	}
}

func loadData(transformedData []TransformedData) {
	// Replace this with your actual data loading logic
	for _, data := range transformedData {
		fmt.Printf("Loading data: %v\n", data)
		// Load data into the destination database or file
	}
}

func main() {
	batchSize := 1000 // Adjust batch size as needed
	batches := extractData(batchSize)
	transformedData := make([]TransformedData, 0)
	var wg sync.WaitGroup

	ctx := context.Background()

	for _, batch := range batches {
		wg.Add(1)
		go transformData(ctx, batch, &wg, &transformedData)
	}

	wg.Wait()
	loadData(transformedData)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
