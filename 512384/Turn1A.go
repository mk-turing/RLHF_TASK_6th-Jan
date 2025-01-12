package main

import (
	"fmt"
)

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
}

// List of voters
var voters []Voter

// List of votes
var votes []Vote

func verifyVoter(voterID string) (bool, error) {
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
	//if verificationResult["is_verified"].(bool) {
	//	return true, nil
	//}
	//
	//return false, fmt.Errorf("Voter not verified")
	return true, nil
}

func registerVoter(voterID string, name string, age int) error {
	isVerified, err := verifyVoter(voterID)
	if err != nil {
		return err
	}
	if !isVerified {
		return fmt.Errorf("Voter ID %s is not verified", voterID)
	}

	newVoter := Voter{ID: voterID, Name: name, Age: age, Registered: true}
	voters = append(voters, newVoter)
	fmt.Println("Voter registered successfully.")
	return nil
}

func castVote(voterID string, candidate string) error {
	for _, voter := range voters {
		if voter.ID == voterID && voter.Registered {
			newVote := Vote{VoterID: voterID, Candidate: candidate}
			votes = append(votes, newVote)
			fmt.Println("Vote recorded successfully.")
			return nil
		}
	}
	return fmt.Errorf("Voter not found or not registered")
}

func main() {
	err := registerVoter("12345", "John Doe", 30)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = castVote("12345", "Candidate A")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Voters:", voters)
	fmt.Println("Votes:", votes)
}
