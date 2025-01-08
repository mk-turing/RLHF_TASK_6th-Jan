package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/compute/v1"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"google.golang.org/api/option"
)

type cloudProvider struct {
	name   string
	client interface{}
}

var providers = []cloudProvider{
	{
		name:   "AWS",
		client: &ec2.Client{
			//Config: config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1")),
		},
	},
	{
		name:   "GCP",
		client: compute.NewService(context.TODO(), option.WithCredentialsFile("path/to/your/gcp-credentials.json")),
	},
	// Add Azure integration here
}

type result struct {
	latency time.Duration
	error   error
}

func main() {
	var provider string
	var concurrency int
	var duration time.Duration
	flag.StringVar(&provider, "p", "AWS", "cloud provider (AWS, GCP)")
	flag.IntVar(&concurrency, "c", 100, "number of concurrent requests")
	flag.DurationVar(&duration, "d", 10*time.Second, "duration of the test")
	flag.Parse()

	fmt.Println("Starting network performance benchmarking on", provider)

	// Select the appropriate cloud provider
	_, ok := getProvider(provider)
	if !ok {
		fmt.Println("Invalid cloud provider specified.")
		return
	}

	// Define the URL to be benchmarked
	url := "https://catfact.ninja/fact"

	// Create a channel to collect results
	results := make(chan result, concurrency)

	// Start the benchmark
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(concurrency)

	// Goroutine to collect the results while requests are being processed
	go func() {
		var totalLatency time.Duration
		var numRequests int
		var numErrors int

		for r := range results {
			if r.error != nil {
				numErrors++
			} else {
				totalLatency += r.latency
				numRequests++
			}
		}

		avgLatency := float64(totalLatency) / float64(numRequests)
		throughput := float64(numRequests) / duration.Seconds()
		reliability := float64(numRequests) / float64(numRequests+numErrors)

		// Print results
		fmt.Println("Benchmark completed on", provider)
		fmt.Printf("Average Latency: %.2f ms\n", avgLatency/float64(time.Millisecond))
		fmt.Printf("Throughput: %.2f requests/second\n", throughput)
		fmt.Printf("Reliability: %.2f%%\n", reliability*100)
	}()

	// Send requests concurrently
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			client := http.DefaultClient
			for startTime.Add(duration).After(time.Now()) {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					results <- result{error: err}
					continue
				}
				start := time.Now()
				resp, err := client.Do(req)
				if err != nil {
					results <- result{error: err}
					continue
				}
				resp.Body.Close()
				latency := time.Since(start)
				results <- result{latency: latency}
			}
		}()
	}

	wg.Wait() // Wait for all goroutines to finish
	close(results)
}

func getProvider(name string) (cloudProvider, bool) {
	for _, p := range providers {
		if p.name == name {
			return p, true
		}
	}
	return cloudProvider{}, false
}
