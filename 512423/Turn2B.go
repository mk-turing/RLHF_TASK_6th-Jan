package main
import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	batchSize = 1000 // Batch size for processing
	numShards  = 4    // Number of map shards
)

type SourceData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type TransformedData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	AgeGroup  string `json:"age_group"`
	IsAdult   bool   `json:"is_adult"`
}

// ShardedMap stores transformed data in multiple maps for scalability
type ShardedMap struct {
	maps []map[int]TransformedData
}

func newShardedMap(numShards int) *ShardedMap {
	m := &ShardedMap{
		maps: make([]map[int]TransformedData, numShards),
	}
	for i := range m.maps {
		m.maps[i] = make(map[int]TransformedData)
	}
	return m
}

func (m *ShardedMap) Get(key int) TransformedData {
	shardIndex := key % numShards
	return m.maps[shardIndex][key]
}

func (m *ShardedMap) Set(key int, value TransformedData) {
	shardIndex := key % numShards
	m.maps[shardIndex][key] = value
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

func transformData(sourceData []SourceData, shardedMap *ShardedMap, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, record := range sourceData {
		ageGroup := "Young"
		isAdult := false
		if record.Age >= 18 {
			ageGroup = "Adult"
			isAdult = true
		}

		shardedMap.Set(record.ID, TransformedData{
			ID:        record.ID,
			Name:      record.Name,
			AgeGroup:  ageGroup,
			IsAdult:   isAdult,
		})
	}
}

func loadData(shardedMap *ShardedMap) {
	// Replace this with your actual data loading logic
	for _, m := range shardedMap.maps {
		for _, data := range m {
			fmt.Printf("Loading data: %v\n", data)
			// Load data into the destination database or file