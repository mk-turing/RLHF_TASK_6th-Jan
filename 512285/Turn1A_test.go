package _12285

import (
	"fmt"
	"testing"
)

type Database interface {
	Query(string, ...interface{}) ([]map[string]interface{}, error)
	Insert(string, ...interface{}) error
	Update(string, ...interface{}) error
	Delete(string, ...interface{}) error
}

type Service struct {
	db Database
}

func NewService(db Database) *Service {
	return &Service{db}
}

func (s *Service) GetUser(id int) (map[string]interface{}, error) { // Changed return type from *map[string]interface{} to map[string]interface{}
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return rows[0], nil // rows[0] is of type map[string]interface{}
}

type MockDatabase struct {
	queryCalls []interface{}
	queryErr   error
	queryRes   []map[string]interface{}
}

func (m *MockDatabase) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	m.queryCalls = append(m.queryCalls, args)
	return m.queryRes, m.queryErr
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

func TestService_GetUser_Success(t *testing.T) {
	// Arrange
	mockDb := &MockDatabase{
		queryRes: []map[string]interface{}{
			{"id": 1, "name": "Alice"},
		},
	}
	s := NewService(mockDb)

	// Act
	user, err := s.GetUser(1)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if user == nil {
		t.Fatal("expected user to not be nil")
	}

	if user["id"] != 1 || user["name"] != "Alice" {
		t.Fatalf("expected user to be {id: 1, name: Alice}, got: %#v", user)
	}

	if len(mockDb.queryCalls) != 1 {
		t.Fatalf("expected 1 query call, got: %d", len(mockDb.queryCalls))
	}

	args := mockDb.queryCalls[0].([]interface{})
	if len(args) != 1 {
		t.Fatalf("expected 1 query argument, got: %d", len(args))
	}

	if args[0] != 1 {
		t.Fatalf("expected query argument to be 1, got: %v", args[0])
	}
}

func TestService_GetUser_NotFound(t *testing.T) {
	// Arrange
	mockDb := &MockDatabase{
		queryRes: []map[string]interface{}{}, // Simulate no rows found
	}
	s := NewService(mockDb)

	// Act
	_, err := s.GetUser(1)

	// Assert
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "user not found" {
		t.Fatalf("expected error message 'user not found', got: %v", err)
	}
}
