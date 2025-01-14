package config

type User struct {
	ID       string
	Balances []Balance
}

type Balance struct {
	Currency string
	Amount   float64
}

var InitialUsers = []User{
	{
		ID: "user1",
		Balances: []Balance{
			{Currency: "BTC", Amount: 0.5},
			{Currency: "ETH", Amount: 5.0},
		},
	},
	{
		ID: "user2",
		Balances: []Balance{
			{Currency: "BTC", Amount: 1.0},
		},
	},
}
