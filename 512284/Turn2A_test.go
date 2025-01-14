package _12284

import (
	"fmt"
	"strings"
	"testing"
)

type User struct {
	ID    int
	Name  string
	Email string
}

type UserService struct{}

func (u UserService) Save(user User) error {
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if len(user.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}
	return nil
}

func TestUserService_ShouldSaveNewUserSuccessfully(t *testing.T) {
	user := User{Name: "Example User", Email: "example@example.com"}
	userService := UserService{}

	err := userService.Save(user)

	if err != nil {
		t.Errorf("Failed to save the user: %v", err)
	}
}

func TestUserService_ShouldFailToSaveUserWithEmptyEmail(t *testing.T) {
	user := User{Name: "Example User", Email: ""}
	userService := UserService{}

	err := userService.Save(user)

	if err == nil {
		t.Error("Expected an error for empty email")
	}
	if !strings.Contains(err.Error(), "email is required") {
		t.Errorf("Expected error message 'email is required', got %v", err.Error())
	}
}

func TestUserService_ShouldFailToSaveUserWithShortName(t *testing.T) {
	user := User{Name: "A", Email: "example@example.com"}
	userService := UserService{}

	err := userService.Save(user)

	if err == nil {
		t.Error("Expected an error for short name")
	}
	if !strings.Contains(err.Error(), "name must be at least 2 characters long") {
		t.Errorf("Expected error message 'name must be at least 2 characters long', got %v", err.Error())
	}
}

func TestUserService_ShouldSaveUserWithBoundaryNameLength(t *testing.T) {
	user := User{Name: "aa", Email: "example@example.com"}
	userService := UserService{}

	err := userService.Save(user)

	if err != nil {
		t.Errorf("Failed to save user with boundary name length: %v", err)
	}
}
