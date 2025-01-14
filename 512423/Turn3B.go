package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type SourceData struct {
	ID   int
	Name string
	Age  int
}

type TransformedData struct {
	ID        int
	Name      string
	AgeGroup  string
	IsAdult   bool
	Version   int // Version number for conflict resolution
}

const (
	redisAddress  = "redis-cluster:6379" // Replace with Redis cluster address
	redisPassword = ""                    // Replace with Redis cluster password (if set)
	redisDB       = 0                    // Replace with Redis cluster database (if set)
	syncInterval  = 10 * time.Second    // Interval for synchronizing local maps with Redis
)

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

func transformData(ctx context.Context, sourceData []SourceData, wg *sync.WaitGroup, localMap *map[int]TransformedData, redisClient *redis.Client) {
	defer wg.Done()

	for _, record := range sourceData {
		ageGroup := "Young"
		isAdult := false
		if record.Age >= 18 {
			ageGroup = "Adult"
			isAdult = true
		}

		(*localMap)[record.ID] = TransformedData{
			ID:        record.ID,
			Name:      record.Name,
			AgeGroup:  ageGroup,
			IsAdult:   isAdult,
			Version:   1, // Initial version
		}
	}

	// Start a goroutine to synchronize local map with Redis periodically
	go syncLocalMapWithRedis(ctx, localMap, redisClient)
}

func syncLocalMapWithRedis(ctx context.Context, localMap *map[int]TransformedData, redisClient *redis.Client) {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for key, data := range *localMap {
				// Get the value from Redis
				redisData, err := redisClient.Get(ctx, fmt.Sprintf("%d", key)).Result()
				if err == redis.Nil {
					// Key doesn't exist in Redis, set it
					if err := redisClient.Set(ctx, fmt.Sprintf("%d", key), data, 0).Err(); err != nil {
						log.Printf("Error setting key %d in Redis: %v", key, err)
					}
				} else if err != nil {
					log.Printf("Error getting key %d from Redis: %v", key, err)
				} else {
					// Key exists in Redis, check version and resolve conflicts
					redisTransformedData, err := decodeTransformedData(redisData)
					if err != nil {
						log.Printf("Error decoding Redis data for key %d: %v", key, err)
						continue
					}

					if redisTransformedData.Version > data.Version {
						// Redis value is newer, update local map
						(*localMap)[key] = redisTransformedData
					} else if redisTransformedData.Version < data.Version {
						// Local value is newer, update Redis
						if err := redisClient.Set(ctx, fmt.Sprintf("%d", key), data, 0).Err(); err != nil {
							log.Printf("Error setting key %d in Redis: %v", key, err)
						}
					} else {
						// Version numbers are the same, resolve conflict using custom logic
						// For this example, we'll just take the latest value
						(*localMap)[key] = data
					}
				}
			}
		}
	}
}

func loadData(transformedData map[int]TransformedData) {
	// Replace this with your actual data loading logic
	for _, data := range transformedData {
		fmt.Printf("Loading data: %v\n", data)
		// Load data into the destination database or file
	}
}

func decodeTransformedData(redisData string) (TransformedData, error) {
	// Implement your decoding logic here
	// For this example, we'll just assume the data is encoded in JSON
	var data TransformedData
	//err := json.Unmarshal([]byte(redisData), &data)
	return data, nil // Return nil for simplicity
}

func main() {
	batchSize := 1000 // Adjust batch size as needed
	batches := extractData(batchSize)
	localMap := make(map[int]TransformedData)
	var wg sync.WaitGroup

	ctx := context.Background()

	// Create a Redis client
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{redisAddress},
		Password: redisPassword,
		DB:       redisDB,
	})
	defer redisClient.Close()

	for _, batch := range batches {
		wg.Add(1)
		go transformData(ctx, batch, &wg, &localMap, redisClient)
	}