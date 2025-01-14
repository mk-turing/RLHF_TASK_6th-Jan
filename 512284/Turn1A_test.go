package _12284

import "testing"

func Test_Account_CreateAccount_ValidData(t *testing.T) {
	account := account.NewAccount("John Doe", "johndoe@example.com", "123456")
	err := account.CreateAccount()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
