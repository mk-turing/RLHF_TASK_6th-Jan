package transaction_log

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Transaction struct {
	ID             string
	Operation      string // e.g., "placeOrder", "updateBalance"
	UserID         string
	OrderID        string
	OrderDetails   *OrderDetails
	BalanceDetails *BalanceDetails
	Timestamp      string
}

type OrderDetails struct {
	Cryptocurrency string
	Type           string // "buy" or "sell"
	Price          float64
	Amount         float64
}

type BalanceDetails struct {
	Cryptocurrency string
	Amount         float64
}

var (
	transactionLog      = make(map[string]*Transaction, 0)
	transactionLogMutex sync.RWMutex
	logFilePath         = "transaction_log.json"
)

func SaveTransaction(transaction *Transaction) error {
	transactionLogMutex.Lock()
	defer transactionLogMutex.Unlock()

	transactionLog[transaction.ID] = transaction

	bytes, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("error marshalling transaction: %v", err)
	}

	err = ioutil.WriteFile(logFilePath, append(bytes, '\n'), 0644)
	if err != nil {
		return fmt.Errorf("error writing transaction to file: %v", err)
	}

	return nil
}

func ReplayLogs() error {
	file, err := os.Open(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Log file does not exist, no need to replay
		}
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer file.Close()

	scanner := json.NewDecoder(file)
	for {
		var transaction Transaction
		err := scanner.Decode(&transaction)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("error decoding transaction: %v", err)
		}

		// Here you would process and replay each transaction.
		// For demonstration, we'll just log it.
		log.Printf("Replaying transaction: %+v\n", transaction)
	}

	return nil
}
