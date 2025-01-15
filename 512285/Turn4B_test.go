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

type Notify interface {
	SendNotification(to string, message string) error
}

type Service struct {
	db     Database
	cache  Cache
	notify Notify
}

func NewService(db Database, cache Cache, notify Notify) *Service {
	return &Service{db, cache, notify}
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

func (s *Service) SomeBusinessLogic() error {
	if err := s.notify.SendNotification("alice@example.com", "Some notification"); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
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

type MockNotify struct {
	sendCalls []interface{}
	sendErr   error
}

func (m *MockNotify) SendNotification(to string, message string) error {
	m.sendCalls = append(m.sendCalls, []interface{}{to, message})
	return m.sendErr
}

func TestService_SomeBusinessLogic_SendsNotification(t *testing.T) {
	// Arrange
	mockDb := &MockDatabase{}
	mockCache := &MockCache{}
	mockNotify := &MockNotify{}
	s := NewService(mockDb, mockCache, mockNotify)

	// Act
	if err := s.SomeBusinessLogic(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Assert
	if len(mockNotify.sendCalls) != 1 {
		t.Fatalf("expected 1 notification send call, got: %d", len(mockNotify.sendCalls))
	}

	sendArgs := mockNotify.sendCalls[0].([]interface{})
	if len(sendArgs) != 2 {
		t.Fatalf("expected 2 notification send arguments, got: %d", len(sendArgs))
	}

	if sendArgs[0] != "alice@example.com" {
		t.Fatalf("expected notification recipient to be 'alice@example.com', got: %v", sendArgs[0])
	}

	if sendArgs[1] != "Some notification" {
		t.Fatalf("expected notification message to be 'Some notification', got: %v", sendArgs[1])
	}
}

func TestService_SomeBusinessLogic_HandlesNotificationFailure(t *testing.T) {
	// Arrange
	mockDb := &MockDatabase{}
	mockCache := &MockCache{}
	mockNotify := &MockNotify{
		sendErr: fmt.Errorf("notification send failed"),
	}
	s := NewService(mockDb, mockCache, mockNotify)

	// Act
	err := s.SomeBusinessLogic()

	// Assert
	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "failed to send notification: notification send failed" {
		t.Fatalf("expected error message 'failed to send notification: notification send failed', got: %v", err)
	}
}
