package users

import (
	"crypto-balance-service/config"
	"strconv"
)

// UserBalances stores user balances as a map
var userBalances = map[string]map[string]float64{}

func GetUserBalance(userID, crypto string) (float64, bool) {
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

	if userBalances[userID] == nil {
		userBalances[userID] = map[string]float64{}
	}

	userBalances[userID][crypto] += value

	return nil
}

func parseBalance(amount string) (float64, error) {
	return strconv.ParseFloat(amount, 64)
}

func Init() {
	userBalances = make(map[string]map[string]float64)
	for _, user := range config.InitialUsers {
		userBalances[user.ID] = make(map[string]float64)
		for _, balance := range user.Balances {
			userBalances[user.ID][balance.Currency] = balance.Amount
		}
	}
}
