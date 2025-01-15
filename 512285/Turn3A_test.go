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

type MockDatabase struct {
	queryCalls   []interface{}
	queryErr     error
	queryRes     []map[string]interface{}
	expectations []*expectation
}

type expectation struct {
	query string
	args  []interface{}
	err   error
	res   []map[string]interface{}
}

func (m *MockDatabase) SetQueryExpectation(query string, args ...interface{}) {
	m.expectations = append(m.expectations, &expectation{
		query: query,
		args:  args,
		err:   m.queryErr,
		res:   m.queryRes,
	})
}

func (m *MockDatabase) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	for _, exp := range m.expectations {
		if exp.query == query && m.equalArgs(exp.args, args) {
			m.queryCalls = append(m.queryCalls, args)
			return exp.res, exp.err
		}
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

func (m *MockDatabase) equalArgs(a, b []interface{}) bool {
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

func TestService_GetUser_CacheHit_WithDynamicMock(t *testing.T) {
	// Arrange
	mockDb := &MockDatabase{}
	mockCache := &MockCache{}
	s := NewService(mockDb, mockCache)

	mockCache.SetGetExpectation("user_1", map[string]interface{}{"id": 1, "name": "Alice"}, nil)

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

	if len(mockCache.getCalls) != 1 {
		t.Fatalf("expected 1 cache get call, got: %d", len(mockCache.getCalls))
	}

	if mockCache.getCalls[0] != "user_1" {
		t.Fatalf("expected cache key to be 'user_1', got: %v", mockCache.getCalls[0])
	}

	if len(mockDb.queryCalls) != 0 {
		t.Fatalf("expected 0 database query calls, got: %d", len(mockDb.queryCalls))
	}
}

func TestService_GetUser_CacheMiss_DatabaseHit_WithDynamicMock(t *testing.T) {
	// Arrange
	mockDb := &MockDatabase{}
	mockCache := &MockCache{}
	s := NewService(mockDb, mockCache)

	// Set cache miss expectation
	mockCache.SetGetExpectation("user_1", nil, fmt.Errorf("cache miss"))

	// Set query result
	mockDb.queryRes = []map[string]interface{}{{"id": 1, "name": "Alice"}}

	// Set the database query expectation for the query to match
	mockDb.SetQueryExpectation("SELECT * FROM users WHERE id = ?", 1)

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

	if len(mockCache.getCalls) != 1 {
		t.Fatalf("expected 1 cache get call, got: %d", len(mockCache.getCalls))
	}

	if mockCache.getCalls[0] != "user_1" {
		t.Fatalf("expected cache key to be 'user_1', got: %v", mockCache.getCalls[0])
	}

	if len(mockDb.queryCalls) != 1 {
		t.Fatalf("expected 1 database query call, got: %d", len(mockDb.queryCalls))
	}

	args := mockDb.queryCalls[0].([]interface{})
	if len(args) != 1 {
		t.Fatalf("expected 1 query argument, got: %d", len(args))
	}

	if args[0] != 1 {
		t.Fatalf("expected query argument to be 1, got: %v", args[0])
	}

	if len(mockCache.setCalls) != 1 {
		t.Fatalf("expected 1 cache set call, got: %d", len(mockCache.setCalls))
	}

	setArgs := mockCache.setCalls[0].([]interface{})
	if len(setArgs) != 3 {
		t.Fatalf("expected 3 cache set arguments, got: %d", len(setArgs))
	}

	if setArgs[0] != "user_1" {
		t.Fatalf("expected cache set key to be 'user_1', got: %v", setArgs[0])
	}
}
