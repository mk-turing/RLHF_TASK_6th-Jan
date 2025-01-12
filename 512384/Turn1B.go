package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Voter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	VoterID    string `json:"voter_id"`
	Verified   bool   `json:"verified"`
	Voted      bool   `json:"voted"`
	LastVoted time.Time
}

type VoterInfo struct {
	VoterID string `json:"voter_id"`
}

type VoterResponse struct {
	Success bool       `json:"success"`
	Data    VoterInfo `json:"data"`
}

func verifyVoter(voterID string) (*VoterInfo, error) {
	// Replace this with the actual API endpoint and token
	apiURL := "https://your-voter-verification-api.com/verify"
	apiToken := "your-api-token"

	// Create a HTTP request
	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Add the API token to the request header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	// Set the content type to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a JSON body with the voter ID
	body := VoterInfo{VoterID: voterID}
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	// Create a HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Parse the JSON response
	var voterResp VoterResponse
	err = json.Unmarshal(data, &voterResp)
	if err != nil {
		return nil, err
	}

	// Return the voter information if verification is successful
	if voterResp.Success {
		return &voterResp.Data, nil
	} else {
		return nil, fmt.Errorf("Voter verification failed")
	}
}

func main() {
	voters := []Voter{}

	// Read voter IDs from a file (replace this with your actual source)
	file, err := os.Open("voter_ids.txt")
	if err != nil {
		fmt.Println("Error reading voter IDs file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		voterID := scanner.Text()
		voterInfo, err := verifyVoter(voterID)
		if err != nil {
			fmt.Println("Error verifying voter", voterID, ":", err)
			continue
		}

		// Create a new voter struct and add it to the slice
		voters = append(voters, Voter{
			Name:       voterInfo.Name,
			VoterID:    voterInfo.VoterID,
			Verified:   true,
			Voted:      false,
			LastVoted:  time.Time{},