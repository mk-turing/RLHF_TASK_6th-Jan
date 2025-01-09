package main

import (
	"fmt"
	"sync"
	"testing"
)

// A simple struct to store dataset processing result
type ProcessResult struct {
	Total   int
	Average float64
	Errors  []error
}

// A function simulating dataset processing
func processDataset(data []int) (Result, error) {
	sum := 0
	count := 0
	var errors []error

	for _, value := range data {
		if value%2 != 0 {
			sum += value
			count++
		} else {
			errors = append(errors, fmt.Errorf("even number found: %d", value))
		}
	}

	average := 0.0
	if count > 0 {
		average = float64(sum) / float64(count)
	}

	return Result{Total: sum, Average: average, Errors: errors}, nil
}

type Result struct {
	Total   int
	Average float64
	Errors  []error
}

// Concurrent dataset processing with error collection
func processDatasetsConcurrently(datasets [][]int) (results []Result, errs []error) {
	const numGoroutines = 10
	var wg sync.WaitGroup
	resultCh := make(chan Result)
	errorCh := make(chan error)

	for _, dataset := range datasets {
		wg.Add(1)
		go func(dataset []int) {
			defer wg.Done()
			result, err := processDataset(dataset)
			resultCh <- result
			if err != nil {
				errorCh <- err
			}
		}(dataset)
	}

	go func() {
		wg.Wait()
		close(resultCh)
		close(errorCh)
	}()

	for result := range resultCh {
		results = append(results, result)
	}
	for err := range errorCh {
		errs = append(errs, err)
	}

	return
}

func BenchmarkProcessDatasets(b *testing.B) {
	b.ResetTimer()
	b.SetParallelism(10)

	for i := 0; i < b.N; i++ {
		// Generate some test datasets
		datasets := make([][]int, 100)
		for j := 0; j < 100; j++ {
			datasets[j] = generateRandomDataset(10000) // 10,000 elements per dataset
		}

		processDatasetsConcurrently(datasets)
	}
}

// Helper function to generate a random dataset
func generateRandomDataset(size int) []int {
	return make([]int, size, size)
}
