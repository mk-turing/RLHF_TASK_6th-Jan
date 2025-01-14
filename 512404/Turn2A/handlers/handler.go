package handlers

import (
	"crypto-balance-service/order_book"
	"fmt"
	"net/http"
	"strconv"
)

var orderBook = order_book.NewOrderBook()

func HandleOrder(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/placeOrder":
		placeOrder(w, r)
	case "/getOrderBook":
		getOrderBook(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func placeOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("orderID")
	userID := r.URL.Query().Get("userID")
	crypto := r.URL.Query().Get("crypto")
	orderType := r.URL.Query().Get("type")
	priceStr := r.URL.Query().Get("price")
	amountStr := r.URL.Query().Get("amount")

	if orderID == "" || userID == "" || crypto == "" || orderType == "" || priceStr == "" || amountStr == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order := order_book.Order{
		ID:             orderID,
		UserID:         userID,
		Cryptocurrency: crypto,
		Type:           orderType,
		Price:          price,
		Amount:         amount,
		Status:         "open",
	}

	err = orderBook.PlaceOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order placed successfully"))

	orderBook.MatchOrders()
}

func getOrderBook(w http.ResponseWriter, r *http.Request) {
	buys, sells := orderBook.GetOrders()
	fmt.Fprintf(w, "Buy Orders:\n")
	for _, buy := range buys {
		fmt.Fprintf(w, "ID: %s, Price: %.2f, Amount: %.2f\n", buy.ID, buy.Price, buy.Amount)
	}

	fmt.Fprintf(w, "\nSell Orders:\n")
	for _, sell := range sells {
		fmt.Fprintf(w, "ID: %s, Price: %.2f, Amount: %.2f\n", sell.ID, sell.Price, sell.Amount)
	}
}
