package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// ElectionType represents the type of election (e.g., General, Referendum)
type ElectionType string

const (
	General    ElectionType = "General"
	Referendum ElectionType = "Referendum"
)

// Voter struct to hold voter information
type Voter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Registered bool   `json:"registered"`
}

// Vote struct to hold vote details for each election type
type Vote struct {
	VoterID      string
	ElectionType ElectionType
	Candidates   []string // For General Election, multiple candidates; For Referendum, only one choice
}

// Election struct to represent an election with its votes and other details
type Election struct {
	Type        ElectionType
	Description string
	Candidates  []string
	Votes       []Vote
}

var voters []Voter
var elections []Election

func verifyVoter(voterID string) (bool, error) {
	// Verifier API call remains the same
}

func getElectionData(electionType ElectionType) ([]string, error) {
	url := fmt.Sprintf("https://api.voterverification.com/elections?type=%s", electionType)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var electionData map[string]interface{}
	err = json.Unmarshal(body, &electionData)
	if err != nil {
		return nil, err
	}

	// Handle API response differently based on election type
	switch electionType {
	case General:
		// For General Election, expect an array of candidates
		candidates := electionData["candidates"].([]interface{})
		var candidateNames []string
		for _, candidate := range candidates {
			candidateNames = append(candidateNames, candidate.(string))
		}
		return candidateNames, nil
	case Referendum:
		// For Referendum, expect a single candidate name
		return []string{electionData["question"].(string)}, nil
	default:
		return nil, fmt.Errorf("Invalid election type: %s", electionType)
	}
}

func registerVoter(voterID string, name string, age int) error {
	// Voter registration remains the same
}

func castVote(voterID string, electionType ElectionType, candidates []string) error {
	// Find the matching election based on the election type
	for _, election := range elections {
		if election.Type == electionType {
			// Validate the candidates for the specific election type
			var validCandidates []string
			for _, candidate := range election.Candidates {
				if contains(candidates, candidate) {
					validCandidates = append(validCandidates, candidate)
				}
			}
			if len(validCandidates) == 0 {
				return fmt.Errorf("Invalid candidates for election type %s", electionType)
			}

			newVote := Vote{VoterID: voterID, ElectionType: electionType, Candidates: validCandidates}
			election.Votes = append(election.Votes, newVote)
			fmt.Println("Vote recorded successfully.")
			return nil
		}
	}
	return fmt.Errorf("Election type not found or not registered")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}