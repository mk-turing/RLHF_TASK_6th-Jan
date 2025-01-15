package _12285

import (
	"fmt"
	"testing"
	"time"
)

type Database interface {
	Query(string, ...interface{}) ([]map[string]interface{}, error)
	Insert(string, ...interface{}) error
	Update(string, ...interface{}) error
	Delete(string, ...interface{}) error
}

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

type EmailService interface {
	SendEmail(to string, subject string, body string) error
}

type Service struct {
	db    Database
	cache Cache
	email EmailService
}

func NewService(db Database, cache Cache, email EmailService) *Service {
	return &Service{db, cache, email}
}

func (s *Service) GetUser(id int) (map[string]interface{}, error) {
	// First, try to retrieve user data from the cache
	key := fmt.Sprintf("user_%d", id)
	userData, err := s.cache.Get(key)
	if err == nil {
		return userData.(map[string]interface{}), nil
	}

	// If not in cache, query the database
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	user := rows[0]

	// Store the retrieved user data in the cache
	err = s.cache.Set(key, user, time.Hour)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) NotifyUser(id int) error {
	// For this example, assume GetUser returns necessary details
	user, err := s.GetUser(id)
	if err != nil {
		return err
	}

	to := user["email"].(string)
	subject := "Welcome to Our Service"
	body := "Hi, " + user["name"].(string) + "!\nWelcome to our service.\nBest,\nThe Team"

	return s.email.SendEmail(to, subject, body)
}

type MockDatabase struct {
	queryCalls   []interface{}
	queryErr     error
	queryRes     []map[string]interface{
	"id":    1,
	"name":  "Alice",
	"email": "alice@example.com",
}
}

func (m *MockDatabase) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	if query == "SELECT * FROM users WHERE id = ?" {
		return m.queryRes, m.queryErr
	}
	return nil, fmt.Errorf("unexpected query: %s with args: %v", query, args)
}

func (m *MockDatabase) Insert(query string, args ...interface{}) error { // Added missing method
	return nil
}

func (m *MockDatabase) Update(query string, args ...interface{}) error { // Added missing method
	return nil
}

func (m *MockDatabase) Delete(query string, args ...interface{}) error { // Added missing method
	return nil
}

type expectation struct {
	query string
	args  []interface{}
	err   error
	res   []map[string]interface{}
}

type MockCache struct {
	setCalls     []interface{}
	setErr       error
	getCalls     []interface{}
	getRes       map[string]interface{}
	getErr       error
	deleteCalls  []interface{}
	deleteErr    error
	expectations []*expectation
}

func (m *MockCache) SetGetExpectation(key string, res interface{}, err error) {
	m.expectations = append(m.expectations, &expectation{
		query: "Get",
		args:  []interface{}{key},
		err:   err,
		res:   []map[string]interface{}{{"key": key, "value": res}},
	})
}

func (m *MockCache) Get(key string) (interface{}, error) {
	for _, exp := range m.expectations {
		if exp.query == "Get" && m.equalArgs([]interface{}{key}, exp.args) {
			m.getCalls = append(m.getCalls, key)
			return exp.res[0]["value"], exp.err
		}
	}
	return nil, fmt.Errorf("unexpected get: %v", key)
}

func (m *MockCache) Set(key string, value interface{}, expiration time.Duration) error {
	m.setCalls = append(m.setCalls, []interface{}{key, value, expiration})
	return m.setErr
}

func (m *MockCache) Delete(key string) error {
	m.deleteCalls = append(m.deleteCalls, key)
	return m.deleteErr
}

func (m *MockCache) equalArgs(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type MockEmailService struct {
	sendEmailCalls []interface{}
	sendEmailErr   error
}

func (m *MockEmailService) SendEmail(to string, subject string, body string) error {
	m.sendEmailCalls = append(m.sendEmailCalls, []interface{}{to, subject, body})
	return m.sendEmailErr
}

func (m *MockEmailService) SetSendEmailExpectation(err error) {
	m.sendEmailErr = err
}

func (m *MockEmailService) equalArgs(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestService_NotifyUser_SuccessfulEmail(t *testing.T) {
	// Arrange
	mockDb := &MockDatabase{}
	mockCache := &MockCache{}
	mockEmail := &MockEmailService{}
	s := NewService(mockDb, mockCache, mockEmail)

	// Set user data in the cache
	mockCache.SetGetExpectation("user_1", map[string]interface{}{"id": 1, "name": "Alice", "email": "alice@example.com"}, nil)

	// Set email service to succeed
	mockEmail.SetSendEmailExpectation(nil)

	// Act
	err := s.NotifyUser(1)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(mockEmail.sendEmailCalls) != 1 {
		t.Fatalf("expected 1 email send call, got: %d", len(mockEmail.sendEmailCalls))
	}

	args := mockEmail.sendEmailCalls[0].([]interface{})
	if len(args) != 3 {
		t.Fatalf("expected 3 email send arguments, got: %d", len(args))
	}

	if args[0] != "alice@example.com" {
		t.Fatalf("expected email to be 'alice@example.com', got: %v", args[0])
	}

	if args[1].(string) != "Welcome to Our Service" {
		t.Fatalf("expected subject to be 'Welcome to Our Service', got: %v", args[1])
	}

	if args[2].(string) != "Hi, Alice!\nWelcome to our service.\nBest,\nThe Team" {
		t.Fatalf("expected body to be 'Hi, Alice!\nWelcome to our service.\nBest,\nThe Team', got: %v", args[2])
	}
}

func TestService_NotifyUser_EmailFailure(t *testing.T) {
	// Arrange
	mockDb := &MockDatabase{}
	mockCache := &MockCache{}
	mockEmail := &MockEmailService{}
	s := NewService(mockDb, mockCache, mockEmail)

	// Set user data in the cache
	mockCache.SetGetExpectation("user_1", map[string]interface{}{"id": 1, "name": "Alice", "email": "alice@example.com"}, nil)

	// Set email service to fail
	mockEmail.SetSendEmailExpectation(fmt.Errorf("email service error"))

	// Act
	err := s.NotifyUser(1)

	// Assert
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "email service error" {
		t.Fatalf("expected error message 'email service error', got: %v", err)
	}

	if len(mockEmail.sendEmailCalls) != 1 {
		t.Fatalf("expected 1 email send call, got: %d", len(mockEmail.sendEmailCalls))
	}
}