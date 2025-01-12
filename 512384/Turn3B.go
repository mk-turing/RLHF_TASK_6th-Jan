package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Mutex lock to protect shared data
var mu sync.Mutex

func addVote(election *Election, voterID string, candidate string) {
	mu.Lock()
	defer mu.Unlock()
	newVote := Vote{VoterID: voterID, Candidate: candidate}
	election.Votes = append(election.Votes, newVote)
	fmt.Printf("Vote added for %s: %s\n", candidate, voterID)
}

func tallyVotes(election *Election) {
	var votes map[string]int
	mu.Lock()
	votes = make(map[string]int, len(election.Votes))
	for _, vote := range election.Votes {
		votes[vote.Candidate]++
	}
	mu.Unlock()
	fmt.Println("--- Tally Result ---")
	for candidate, count := range votes {
		fmt.Printf("Candidate: %s - Votes: %d\n", candidate, count)
	}
}

func voterProcess(election *Election, voterID string, candidate string) {
	// Simulate voting delay
	time.Sleep(time.Duration(1000+rand.Intn(1000)) * time.Millisecond)
	addVote(election, voterID, candidate)
}

func main() {
	election := &Election{
		ID:    "general1",
		Type:  "general",
		Title: "2024 General Election",
		Votes: []Vote{},
	}
	candidates := []string{"Candidate A", "Candidate B", "Candidate C"}

	// Start voter processes concurrently
	numVoters := 20
	for i := 1; i <= numVoters; i++ {
		voterID := fmt.Sprintf("voter%d", i)
		candidate := candidates[rand.Intn(len(candidates))]
		go voterProcess(election, voterID, candidate)
	}

	// Wait for all voter processes to complete
	time.Sleep(time.Second * 5) // Add sufficient time for all voters to complete voting

	// Tally votes after all voting is complete
	tallyVotes(election)
}
