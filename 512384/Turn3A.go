package main

import (
	"context"
	"fmt"
	"sync"
)

var mu sync.Mutex     // Mutex to ensure thread safety
var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish

type Election struct {
	ID                string                 `json:"id"`
	Type              string                 `json:"type"` // e.g., "general", "referendum"
	Title             string                 `json:"title"`
	Voters            []Voter                `json:"voters"`
	Votes             []Vote                 `json:"votes"`
	VerificationRules map[string]interface{} `json:"verification_rules"`
}

type Voter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Registered bool   `json:"registered"`
}

type Vote struct {
	VoterID   string
	Candidate string
	Option    string // For referendums
}

var elections []Election

func verifyVoter(election *Election, voterID string) (bool, error) {
	return true, nil
}

func registerVoterConcurrently(ctx context.Context, electionID string, voterID string, name string, age int) error {
	mu.Lock()
	defer mu.Unlock()

	for i, election := range elections {
		if election.ID == electionID {
			isVerified, err := verifyVoter(&election, voterID)
			if err != nil {
				return err
			}
			if !isVerified {
				return fmt.Errorf("Voter ID %s is not verified", voterID)
			}

			newVoter := Voter{ID: voterID, Name: name, Age: age, Registered: true}
			election.Voters = append(election.Voters, newVoter)
			elections[i] = election
			fmt.Printf("Voter registered successfully for election %s.\n", election.Title)
			return nil
		}
	}
	return fmt.Errorf("Election with ID %s not found", electionID)
}

func castVoteConcurrently(ctx context.Context, electionID string, voterID string, candidate string, option string) error {
	mu.Lock()
	defer mu.Unlock()

	for i, election := range elections {
		if election.ID == electionID {
			for _, voter := range election.Voters {
				if voter.ID == voterID && voter.Registered {
					newVote := Vote{VoterID: voterID, Candidate: candidate, Option: option}
					election.Votes = append(election.Votes, newVote)
					elections[i] = election
					fmt.Printf("Vote recorded successfully for election %s.\n", election.Title)
					return nil
				}
			}
			return fmt.Errorf("Voter not found or not registered")
		}
	}
	return fmt.Errorf("Election with ID %s not found", electionID)
}

func tallyVotes(ctx context.Context) {
	mu.Lock()
	defer mu.Unlock()

	for _, election := range elections {
		voteCounts := make(map[string]int)
		for _, vote := range election.Votes {
			if election.Type == "general" {
				voteCounts[vote.Candidate]++
			} else if election.Type == "referendum" {
				voteCounts[vote.Option]++
			}
		}

		fmt.Printf("Tally for Election %s:\n", election.Title)
		for candidate, count := range voteCounts {
			fmt.Printf("%s: %d\n", candidate, count)
		}
		fmt.Println("---------------------")
	}
}

func createElection(electionID string, electionType string, title string, verificationRules map[string]interface{}) {
	newElection := Election{
		ID:                electionID,
		Type:              electionType,
		Title:             title,
		Voters:            []Voter{},
		Votes:             []Vote{},
		VerificationRules: verificationRules,
	}
	elections = append(elections, newElection)
	fmt.Println("Election created successfully.")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create elections
	createElection("general1", "general", "2024 General Election", map[string]interface{}{})
	createElection("referendum1", "referendum", "2024 Referendum", map[string]interface{}{})

	// Register voters concurrently
	wg.Add(2) // Add 2 tasks for registering voters
	go func() {
		defer wg.Done()
		if err := registerVoterConcurrently(ctx, "general1", "12345", "John Doe", 30); err != nil {
			fmt.Println("Error registering voter:", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := registerVoterConcurrently(ctx, "referendum1", "67890", "Jane Smith", 18); err != nil {
			fmt.Println("Error registering voter:", err)
		}
	}()
	wg.Wait()

	// Cast votes concurrently
	wg.Add(2) // Add 2 tasks for casting votes
	go func() {
		defer wg.Done()
		if err := castVoteConcurrently(ctx, "general1", "12345", "Candidate A", ""); err != nil {
			fmt.Println("Error casting vote:", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := castVoteConcurrently(ctx, "referendum1", "67890", "", "Option Yes"); err != nil {
			fmt.Println("Error casting vote:", err)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Tally the votes after all operations are complete
	tallyVotes(ctx)
}
