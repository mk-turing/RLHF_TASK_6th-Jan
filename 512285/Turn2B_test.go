package _12285

import (
	"errors"
	"testing"
)

type UserService struct {
	UserDB      UserDB
	AuthService AuthService
}

func (s *UserService) CreateUser(email, password string) error {
	if _, err := s.UserDB.GetUserByEmail(email); err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := s.AuthService.HashPassword(password)
	if err != nil {
		return err
	}

	user := &User{Email: email, Password: hashedPassword}
	if err := s.UserDB.CreateUser(user); err != nil {
		return err
	}

	return nil
}

type User struct {
	ID       int
	Email    string
	Password string
}

type UserDB interface {
	CreateUser(*User) error
	GetUserByEmail(string) (*User, error)
}

type MockUserDB struct {
	CreatedUser *User
	UsersByEmail map[string]*User
}

func NewMockUserDB() *MockUserDB {
	return &MockUserDB{
		UsersByEmail: map[string]*User{},
	}
}

func (m *MockUserDB) CreateUser(user *User) error {
	m.CreatedUser = user
	return nil
}

func (m *MockUserDB) GetUserByEmail(email string) (*User, error) {
	user, ok := m.UsersByEmail[email]
	if !ok {
		return nil, nil
	}
	return user, nil
}

type AuthService interface {
	HashPassword(string) (string, error)
}

type MockAuthService struct {
	HashedPasswords map[string]string
	HashPasswordErr error
}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{
		HashedPasswords: map[string]string{},
	}
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	if m.HashPasswordErr != nil {
		return "", m.HashPasswordErr
	}
	hashedPassword, ok := m.HashedPasswords[password]
	if !ok {
		hashedPassword = password // In a real scenario, you would use a proper hashing function
	}
	return hashedPassword, nil
}

func TestUserService_CreateUser(t *testing.T) {
	mockUserDB := NewMockUserDB()
	mockAuthService := NewMockAuthService()
	s := &UserService{
		UserDB:      mockUserDB,
		AuthService: mockAuthService,
	}

	// Success scenario
	email := "test@example.com"
	password := "password"

	if err := s.CreateUser(email, password); err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	createdUser := mockUserDB.CreatedUser
	if createdUser == nil {
		t.Fatal("Expected user to be created")
	}

	if createdUser.Email != email {
		t.Fatalf("Expected user email to be %s, got: %s", email, createdUser.Email)
	}

	hashedPassword, _ := mockAuthService.HashPassword(password)
	if createdUser.Password != hashedPassword {
		t.Fatalf("Expected user password to be hashed, got: %s", createdUser.Password)
	}

	// User already exists scenario
	if err := s.CreateUser(email, password); err == nil {
		t.Fatal("Expected error when user already exists")
	}

	// Hashing password error scenario
	mockAuthService.HashPasswordErr = errors.New("hashing password failed")
	if err := s.CreateUser(email, password); err == nil {
		t.Fatal("Expected error when hashing password fails")
	}