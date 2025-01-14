package _12284

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

type User struct {
	Name  string
	Email string
}

// Example interface for the database service
type DatabaseService interface {
	SaveUser(ctx context.Context, user User) error
}

// Example interface for the email service
type EmailService interface {
	SendWelcomeEmail(ctx context.Context, email string) error
}

type UserService struct {
	db    DatabaseService
	email EmailService
}

func (u UserService) Save(ctx context.Context, user User) error {
	if err := u.db.SaveUser(ctx, user); err != nil {
		return err
	}
	return u.email.SendWelcomeEmail(ctx, user.Email)
}

// Example database service implementation
type MockDatabaseService struct {
	err error
}

func (m MockDatabaseService) SaveUser(ctx context.Context, user User) error {
	return m.err
}

// Example email service implementation
type MockEmailService struct {
	err error
}

func (m MockEmailService) SendWelcomeEmail(ctx context.Context, email string) error {
	return m.err
}

func TestUserService_ShouldSaveAndEmailUser(t *testing.T) {
	user := User{Name: "Example User", Email: "example@example.com"}
	db := MockDatabaseService{}
	email := MockEmailService{}
	userService := UserService{db: db, email: email}

	ctx := context.Background()

	if err := userService.Save(ctx, user); err != nil {
		t.Errorf("Failed to save and email user: %v", err)
	}
}

func TestUserService_ShouldNotSaveWhenDatabaseFails(t *testing.T) {
	user := User{Name: "Example User", Email: "example@example.com"}
	db := MockDatabaseService{err: fmt.Errorf("database error")}
	email := MockEmailService{}
	userService := UserService{db: db, email: email}

	ctx := context.Background()

	if err := userService.Save(ctx, user); err == nil {
		t.Error("Expected an error when database save fails")
	} else if !strings.Contains(err.Error(), "database error") {
		t.Errorf("Expected error message 'database error', got %v", err.Error())
	}
}

func TestUserService_ShouldNotEmailWhenDatabaseFails(t *testing.T) {
	user := User{Name: "Example User", Email: "example@example.com"}
	db := MockDatabaseService{err: fmt.Errorf("database error")}
	email := MockEmailService{}
	userService := UserService{db: db, email: email}

	ctx := context.Background()

	// Check if SaveUser fails as expected
	if err := userService.Save(ctx, user); err == nil {
		t.Error("Expected an error when database save fails")
	} else if !strings.Contains(err.Error(), "database error") {
		t.Errorf("Expected error message 'database error', got %v", err.Error())
	}

	// We should verify that the email service is not called due to the failure in database save
	// You can modify the MockEmailService to track the number of calls to SendWelcomeEmail, for instance
}
