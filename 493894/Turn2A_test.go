package _93894

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

// Component1 interface
type Component1 interface {
	ProcessData([]int) int
}

// RealComponent1 is the real implementation of Component1.
type RealComponent1 struct {
}

func (r *RealComponent1) ProcessData(data []int) int {
	sum := 0
	for _, num := range data {
		sum += num
	}
	return sum
}

// MockComponent1 is a mock implementation of Component1 using testify/mock.
type MockComponent1 struct {
	mock.Mock
}

func (m *MockComponent1) ProcessData(data []int) int {
	args := m.Called(data)
	return args.Int(0)
}

// Component2 interface
type Component2 interface {
	GenerateReport(string) string
}

// RealComponent2 is the real implementation of Component2.
type RealComponent2 struct {
}

func (r *RealComponent2) GenerateReport(data string) string {
	return "Generated Report: " + data
}

// MockComponent2 is a mock implementation of Component2 using testify/mock.
type MockComponent2 struct {
	mock.Mock
}

func (m *MockComponent2) GenerateReport(data string) string {
	args := m.Called(data)
	return args.String(0)
}

func BenchmarkRealComponent1(b *testing.B) {
	svc := &RealComponent1{}
	data := []int{1, 2, 3, 4, 5}
	for i := 0; i < b.N; i++ {
		svc.ProcessData(data)
	}
}

func BenchmarkMockComponent1_Stub(b *testing.B) {
	svc := &MockComponent1{}
	data := []int{1, 2, 3, 4, 5}
	svc.On("ProcessData", data).Return(15)
	for i := 0; i < b.N; i++ {
		svc.ProcessData(data)
	}
}

func BenchmarkMockComponent1_Expect(b *testing.B) {
	svc := &MockComponent1{}
	data := []int{1, 2, 3, 4, 5}
	svc.On("ProcessData", data).Return(15)
	for i := 0; i < b.N; i++ {
		svc.ProcessData(data)
	}
	svc.AssertExpectations(b)
}

func BenchmarkRealComponent2(b *testing.B) {
	svc := &RealComponent2{}
	data := "test data"
	for i := 0; i < b.N; i++ {
		svc.GenerateReport(data)
	}
}

func BenchmarkMockComponent2_Stub(b *testing.B) {
	svc := &MockComponent2{}
	data := "test data"
	svc.On("GenerateReport", data).Return("Generated Report: test data")
	for i := 0; i < b.N; i++ {
		svc.GenerateReport(data)
	}
}

func BenchmarkMockComponent2_Expect(b *testing.B) {
	svc := &MockComponent2{}
	data := "test data"
	svc.On("GenerateReport", data).Return("Generated Report: test data")
	for i := 0; i < b.N; i++ {
		svc.GenerateReport(data)
	}
	svc.AssertExpectations(b)
}
