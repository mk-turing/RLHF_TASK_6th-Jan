package main

import (
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

func extractData() []SourceData {
	// Replace this with your actual data extraction logic
	return []SourceData{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 18},
		// Add more data as needed
	}
}

func transformData(sourceData []SourceData) map[int]TransformedData {
	transformedDataMap := make(map[int]TransformedData)

	var wg sync.WaitGroup
	wg.Add(len(sourceData))

	for _, record := range sourceData {
		go func(r SourceData) {
			defer wg.Done()

			ageGroup := "Young"
			isAdult := false
			if r.Age >= 18 {
				ageGroup = "Adult"
				isAdult = true
			}

			transformedDataMap[r.ID] = TransformedData{
				ID:       r.ID,
				Name:     r.Name,
				AgeGroup: ageGroup,
				IsAdult:  isAdult,
			}
		}(record)
	}

	wg.Wait()
	return transformedDataMap
}

func loadData(transformedDataMap map[int]TransformedData) {
	// Replace this with your actual data loading logic
	for _, data := range transformedDataMap {
		fmt.Printf("Loading data: %v\n", data)
		// Load data into the destination database or file
	}
}

func main() {
	sourceData := extractData()
	transformedDataMap := transformData(sourceData)
	loadData(transformedDataMap)
}
