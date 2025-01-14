package handlers

import (
	"fmt"
	"net/http"

	"crypto-balance-service/users"
)

// HandleRequest handles HTTP requests
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	switch path {
	case "/getBalance":
		getBalance(w, r)
	case "/updateBalance":
		updateBalance(w, r)
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	crypto := r.URL.Query().Get("crypto")

	if userID == "" || crypto == "" {
		http.Error(w, "Missing userID or crypto", http.StatusBadRequest)
		return
	}

	balance, found := users.GetUserBalance(userID, crypto)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Balance not found"))
		return
	}

	fmt.Fprintf(w, "Balance: %f\n", balance)
}

func updateBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	crypto := r.URL.Query().Get("crypto")
	amount := r.URL.Query().Get("amount")

	if userID == "" || crypto == "" || amount == "" {
		http.Error(w, "Missing userID, crypto, or amount", http.StatusBadRequest)
		return
	}

	err := users.UpdateUserBalance(userID, crypto, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Balance updated successfully"))
}
