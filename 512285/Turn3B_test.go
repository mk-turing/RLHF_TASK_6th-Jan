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

type Service struct {
	db    Database
	cache Cache
}

func NewService(db Database, cache Cache) *Service {
	return &Service{db, cache}
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

type MockDatabaseConfig struct {
	QueryError  error
	QueryResult []map[string]interface{}
}

type MockDatabase struct {
	config     *MockDatabaseConfig
	queryCalls []interface{}
}

func NewMockDatabase(config *MockDatabaseConfig) *MockDatabase {
	return &MockDatabase{
		config: config,
	}
}

func (m *MockDatabase) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	m.queryCalls = append(m.queryCalls, args)
	if m.config.QueryError != nil {
		return nil, m.config.QueryError
	}
	return m.config.QueryResult, nil
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

type MockCacheConfig struct {
	GetError  error
	GetResult map[string]interface{}
	SetError  error
}

type MockCache struct {
	config      *MockCacheConfig
	setCalls    []interface{}
	getCalls    []interface{}
	deleteCalls []interface{}
}

func NewMockCache(config *MockCacheConfig) *MockCache {
	return &MockCache{
		config: config,
	}
}

func (m *MockCache) Set(key string, value interface{}, expiration time.Duration) error {
	m.setCalls = append(m.setCalls, []interface{}{key, value, expiration})
	return m.config.SetError
}

func (m *MockCache) Get(key string) (interface{}, error) {
	m.getCalls = append(m.getCalls, key)
	return m.config.GetResult, m.config.GetError
}

func (m *MockCache) Delete(key string) error {
	m.deleteCalls = append(m.deleteCalls, key)
	return nil
}

func TestService_GetUser_DatabaseErrorBasedOnInput(t *testing.T) {
	// Arrange
	mockDb := NewMockDatabase(&MockDatabaseConfig{
		QueryResult: []map[string]interface{}{
			{"id": 1, "name": "Alice"},
		},
	})
	mockCache := NewMockCache(&MockCacheConfig{
		GetError: fmt.Errorf("cache miss"),
	})
	s := NewService(mockDb, mockCache)

	// Act
	// Call GetUser with different input to trigger different behavior
	_, err1 := s.GetUser(1)
	_, err2 := s.GetUser(2)

	// Assert
	if err1 != nil {
		t.Fatalf("expected no error for user 1, got: %v", err1)
	}
	if err2 == nil {
		t.Fatal("expected error for user 2")
	}
}
