package users

import (
	"crypto-balance-service/config"
	"crypto-balance-service/transaction_log"
	"strconv"
	"sync"
	"time"
)

var (
	userBalances     = map[string]map[string]float64{}
	userBalanceMutex sync.Mutex
)

func GetUserBalance(userID, crypto string) (float64, bool) {
	userBalanceMutex.Lock()
	defer userBalanceMutex.Unlock()

	user, ok := userBalances[userID]
	if !ok {
		return 0, false // User not found
	}
	return user[crypto], true
}

func UpdateUserBalance(userID, crypto string, amount string) error {
	value, err := parseBalance(amount)
	if err != nil {
		return err
	}

	userBalanceMutex.Lock()
	defer userBalanceMutex.Unlock()

	if userBalances[userID] == nil {
		userBalances[userID] = map[string]float64{}
	}

	userBalances[userID][crypto] += value

	transaction_log.SaveTransaction(&transaction_log.Transaction{
		ID:             "balanceUpdate" + time.Now().Format("20060102150405"),
		Operation:      "updateBalance",
		UserID:         userID,
		OrderID:        "", // For balance updates, OrderID is empty
		BalanceDetails: &transaction_log.BalanceDetails{crypto, value},
		Timestamp:      time.Now().Format(time.RFC3339),
	})

	return nil
}

func Init() {
	userBalanceMutex.Lock()
	defer userBalanceMutex.Unlock()

	userBalances = make(map[string]map[string]float64)
	for _, user := range config.InitialUsers {
		userBalances[user.ID] = make(map[string]float64)
		for _, balance := range user.Balances {
			userBalances[user.ID][balance.Currency] = balance.Amount
		}
	}
}

func parseBalance(amount string) (float64, error) {
	return strconv.ParseFloat(amount, 64)
}
