package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func stressTest(ctx context.Context, url string, concurrency int, duration time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Make a request to your endpoint
					if err := makeRequest(url); err != nil {
						// Handle or log errors
						log.Printf("Error: %v\n", err)
					}
					// Add delay between requests
					time.Sleep(time.Millisecond * 100)
				}
			}
		}()
	}

	select {
	case <-ctx.Done():
		log.Printf("Test stopped after %v\n", time.Since(start))
	case <-time.After(duration):
		log.Printf("Test completed after %v\n", duration)
	}
}

func makeRequest(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %v", resp.StatusCode)
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	stressTest(ctx, "http://localhost:8080/your-endpoint", 100, 3*time.Minute)
}
