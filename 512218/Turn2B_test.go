package _12218

import (
	"fmt"
	"net/http"
	"strings"
)

// Global vars for test convenience (should be fixed later)
var whiteListDomains = []string{"localhost:8080"}
var apiKeyHeader = "X-Api-Key"
var goodApiKey = "secretKey"

// Sample implementation of parsing logic (should be adjusted based on your app)
func parseQueryParameters(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	param1 := values.Get("param1")
	param2 := values.Get("param2")

	// Perform validations and sanitizations
	param1 = strings.TrimSpace(param1)
	if !isValidParameter(param1) {
		http.Error(w, "Invalid parameter value", http.StatusBadRequest)
		return
	}
	// Business logic
	result := fmt.Sprintf("Parsed param1: %s, param2: %s", param1, param2)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, result)
}

// Updated handler with security checks and sanitizations
func secureQueryParamHandler(w http.ResponseWriter, r *http.Request) {
	// Secure Header Check (always HTTPS should be configured elsewhere)
	if r.Header.Get("X-Forwarded-Proto") != "https" {
		http.Error(w, "HTTPS required", http.StatusUpgradeRequired)
		return
	}

	// Request Origin Check
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = r.Header.Get("Referer")
	}
	if !isAllowedOrigin(origin) {
		http.Error(w, "Origin not allowed", http.StatusForbidden)
		return
	}

	// API Key Check
	apiKey := r.Header.Get(apiKeyHeader)
	if !isValidApiKey(apiKey) {
		http.Error(w, "Invalid API Key", http.StatusUnauthorized)
		return
	}

	parseQueryParameters(w, r)
}

// .... (Previous test cases remained the same)
