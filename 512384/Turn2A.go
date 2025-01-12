package main

import (
	"fmt"
)

// Election struct to hold election metadata
type Election struct {
	ID                string                 `json:"id"`
	Type              string                 `json:"type"` // e.g., "general", "referendum"
	Title             string                 `json:"title"`
	Voters            []Voter                `json:"voters"`
	Votes             []Vote                 `json:"votes"`
	VerificationRules map[string]interface{} `json:"verification_rules"`
}

// Voter struct to hold voter information
type Voter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Registered bool   `json:"registered"`
}

// Vote struct to hold vote details
type Vote struct {
	VoterID   string
	Candidate string
	Option    string // For referendums
}

// List of elections
var elections []Election

func verifyVoter(election *Election, voterID string) (bool, error) {
	return true, nil
	//url := fmt.Sprintf("https://api.voterverification.com/verify?voter_id=%s", voterID)
	//resp, err := http.Get(url)
	//if err != nil {
	//	return false, err
	//}
	//defer resp.Body.Close()
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	return false, err
	//}
	//
	//var verificationResult map[string]interface{}
	//err = json.Unmarshal(body, &verificationResult)
	//if err != nil {
	//	return false, err
	//}
	//
	//isVerified := verificationResult["is_verified"].(bool)
	//
	//// Apply additional rules based on election type
	//if election.Type == "referendum" {
	//	if age := verificationResult["age"].(int); age < 18 {
	//		isVerified = false
	//	}
	//}
	//
	//return isVerified, nil
}

func registerVoter(electionID string, voterID string, name string, age int) error {
	for _, election := range elections {
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
			fmt.Println("Voter registered successfully for election %s.", election.Title)
			return nil
		}
	}
	return fmt.Errorf("Election with ID %s not found", electionID)
}

func castVote(electionID string, voterID string, candidate string, option string) error {
	for _, election := range elections {
		if election.ID == electionID {
			for _, voter := range election.Voters {
				if voter.ID == voterID && voter.Registered {
					newVote := Vote{VoterID: voterID, Candidate: candidate, Option: option}
					election.Votes = append(election.Votes, newVote)
					fmt.Println("Vote recorded successfully for election %s.", election.Title)
					return nil
				}
			}
			return fmt.Errorf("Voter not found or not registered")
		}
	}
	return fmt.Errorf("Election with ID %s not found", electionID)
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
	createElection("general1", "general", "2024 General Election", map[string]interface{}{})
	createElection("referendum1", "referendum", "2024 Referendum", map[string]interface{}{})

	err := registerVoter("general1", "12345", "John Doe", 30)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = castVote("general1", "12345", "Candidate A", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = registerVoter("referendum1", "67890", "Jane Smith", 17)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, election := range elections {
		fmt.Println("Election:", election.Title)
		fmt.Println("Voters:", election.Voters)
		fmt.Println("Votes:", election.Votes)
		fmt.Println("---------------------")
	}
}
