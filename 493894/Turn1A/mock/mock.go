package mock

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockStorage implements the Storage interface
type MockStorage struct {
	mock.Mock
}

// NewMockStorage creates a new MockStorage instance
func NewMockStorage(t *testing.T) *MockStorage {
	mock := new(MockStorage)
	mock.Object = mock
	return mock
}

// GetData mocks the GetData method of the Storage interface
func (m *MockStorage) GetData() ([]byte, error) {
	args := m.Called()
	ret := args.Get(0)
	return ret.([]byte), ret.(error)
}
