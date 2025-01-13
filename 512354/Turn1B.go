package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func fetchData(url string, callback chan<- string) {
	time.Sleep(time.Second * 2)
	data := fmt.Sprintf("Data from %s", url)
	callback <- data
}

func main() {
	wg := &sync.WaitGroup{}

	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		url := r.URL.Query().Get("url")
		callback := make(chan string, 1)
		fetchData(url, callback)

		select {
		case data := <-callback:
			w.Write([]byte(data))
		case <-time.After(time.Second * 5):
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Timeout fetching data"))
		}
	}

	http.HandleFunc("/", handleRequest)
	log.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
