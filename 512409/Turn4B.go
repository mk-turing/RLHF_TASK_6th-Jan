package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Function with multiple return values
func processData(data []byte) (result1 int, result2 float64, result3 string) {
	// Simulate some processing and return multiple results
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	result1 = len(data)
	result2 = float64(rand.Intn(1000)) / 100
	result3 = fmt.Sprintf("Processed: %d bytes", result1)
	return
}

// Data processor
func dataProcessor(dataCh <-chan []byte, resultCh chan<- Result) {
	for data := range dataCh {
		result1, result2, result3 := processData(data)
		result := Result{
			Result1: result1,
			Result2: result2,
			Result3: result3,
		}
		resultCh <- result
	}
}

// Result aggregator
func resultAggregator(resultCh <-chan Result, aggregateCh chan<- AggregateResult) {
	var aggregate AggregateResult
	for result := range resultCh {
		aggregate.Result1 += result.Result1
		aggregate.Result2 += result.Result2
		aggregate.Result3 = append(aggregate.Result3, result.Result3)
	}
	aggregateCh <- aggregate
}

type Result struct {
	Result1 int
	Result2 float64
	Result3 string
}

type AggregateResult struct {
	Result1 int
	Result2 float64
	Result3 []string
}

func main() {
	// Number of goroutines to use for data processing
	const numGoroutines = 100

	// Channel to receive data
	dataCh := make(chan []byte)

	// Channel to send results from data processors
	resultCh := make(chan Result, numGoroutines)

	// Channel to receive aggregated results
	aggregateCh := make(chan AggregateResult)

	// Start data processors
	var wg sync.WaitGroup
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			dataProcessor(dataCh, resultCh)
		}()
	}

	// Start result aggregator
	go func() {
		resultAggregator(resultCh, aggregateCh)
	}()

	// Generate random data and send it to the data processor
	go func() {
		for {
			data := make([]byte, rand.Intn(1000))
			dataCh <- data
		}
	}()

	// Wait for data processors to finish
	wg.Wait()

	// Close channels to signal the end of data processing
	close(dataCh)
	close(resultCh)

	// Receive and print aggregated results
	aggregateResult := <-aggregateCh
	fmt.Printf("Aggregated Result: Result1: %d, Result2: %.2f, Result3: %q\n", aggregateResult.Result1, aggregateResult.Result2, aggregateResult.Result3)
}
