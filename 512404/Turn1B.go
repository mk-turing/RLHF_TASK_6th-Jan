package main

import (
	"fmt"
	"sync"
)

type Balance struct {
	sync.Mutex
	Currency string
	Amount   float64
}

type UserBalances struct {
	sync.Mutex
	Balances map[string]*Balance
}

func (ub *UserBalances) Init() {
	ub.Balances = make(map[string]*Balance)
}

func (ub *UserBalances) GetBalance(userID string, currency string) (float64, error) {
	ub.Lock()
	defer ub.Unlock()

	userBalance, ok := ub.Balances[userID]
	if !ok {
		return 0, fmt.Errorf("User ID '%s' not found", userID)
	}

	currencyBalance, ok := userBalance.Balances[currency]
	if !ok {
		return 0, fmt.Errorf("Currency '%s' not found for user '%s'", currency, userID)
	}

	return currencyBalance.Amount, nil
}

func (ub *UserBalances) UpdateBalance(userID string, currency string, amount float64) error {
	ub.Lock()
	defer ub.Unlock()

	userBalance, ok := ub.Balances[userID]
	if !ok {
		userBalance = &UserBalances{Balances: make(map[string]*Balance)}
		ub.Balances[userID] = userBalance
	}

	currencyBalance, ok := userBalance.Balances[currency]
	if !ok {
		currencyBalance = &Balance{Currency: currency, Amount: 0}
		userBalance.Balances[currency] = currencyBalance
	}

	currencyBalance.Lock()
	defer currencyBalance.Unlock()

	currencyBalance.Amount += amount

	return nil
}

func main() {
	userBalances := &UserBalances{}
	userBalances.Init()

	err := userBalances.UpdateBalance("user1", "BTC", 1.0)
	if err != nil {
		fmt.Println("Error updating balance:", err)
		return
	}

	balance, err := userBalances.GetBalance("user1", "BTC")
	if err != nil {
		fmt.Println("Error getting balance:", err)
		return
	}

	fmt.Println("User 1 BTC balance:", balance)
}
