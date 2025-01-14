package main

import (
	"crypto-balance-service/handlers"
	"crypto-balance-service/transaction_log"
	"log"
	"net/http"
)

func main() {
	err := transaction_log.ReplayLogs()
	if err != nil {
		log.Fatalf("Error replaying logs: %v\n", err)
	}

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
