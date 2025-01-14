package main

import (
	"crypto-balance-service/handlers"
	"crypto-balance-service/order_book"
	"crypto-balance-service/persistence"
	"crypto-balance-service/transaction_log"
	"crypto-balance-service/users"
	"database/sql"
	"log"
	"net/http"
)

func main() {
	// Initialize database connection
	db, err := sql.Open("sqlite3", "./crypto.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load transaction log from database
	transactionLog, err := persistence.LoadTransactionLog(db)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize order book and user balances
	orderBook := order_book.NewOrderBook()
	userBalances := users.NewUserBalances()

	// Replay transaction log to restore state
	transactionLog.Replay(orderBook, userBalances)

	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handlers.HandleRequest),
	}