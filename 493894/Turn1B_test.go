package _93894

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

type MyService interface {
	DoSomething() string
}

// RealService is the actual implementation of MyService.
type RealService struct {
}

func (s *RealService) DoSomething() string {
	// Simulate some work
	time.Sleep(10 * time.Millisecond)
	return "result"
}

// MockService is a mock implementation of MyService using testify/mock.
type MockService struct {
	mock.Mock
}

func (m *MockService) DoSomething() string {
	args := m.Called()
	return args.String(0)
}

func BenchmarkRealService(b *testing.B) {
	for i := 0; i < b.N; i++ {
		svc := &RealService{}
		svc.DoSomething()
	}
}

func BenchmarkMockService_Stub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		svc := &MockService{}
		svc.On("DoSomething").Return("result")
		svc.DoSomething()
	}
}

func BenchmarkMockService_Expect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		svc := &MockService{}
		// Replace 'Expect' with 'On' method
		svc.On("DoSomething").Return("result")
		svc.DoSomething()
		svc.AssertExpectations(b)
	}
}
