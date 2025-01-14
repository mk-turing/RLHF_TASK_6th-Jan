package main

import (
	"crypto-balance-service/handlers"
	"log"
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handlers.HandleOrder),
	}

	log.Println("Cryptocurrency Balance Service starting on :8080")
	go server.ListenAndServe()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
