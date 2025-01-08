package _93894

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// Mocking definitions for MyDataProcessor and MyDatabaseService
// ... (Same as before)

type MyDataProcessor interface {
	ProcessData(data string) string
}

// RealDataProcessor is the actual implementation of MyDataProcessor
type RealDataProcessor struct {
}

func (p *RealDataProcessor) ProcessData(data string) string {
	// Simulate some work
	time.Sleep(10 * time.Millisecond)
	return data
}

type MockDataProcessor struct {
	mock.Mock
}

func (m *MockDataProcessor) ProcessData(data string) string {
	args := m.Called(data)
	return args.String(0)
}

// Mocking definitions for MyDatabaseService
// ... (Same as before)

type MyDatabaseService interface {
	SaveData(data string) error
}

// RealDatabaseService is the actual implementation of MyDatabaseService
type RealDatabaseService struct {
}

func (s *RealDatabaseService) SaveData(data string) error {
	// Simulate some work
	time.Sleep(15 * time.Millisecond)
	return nil
}

type MockDatabaseService struct {
	mock.Mock
}

func (m *MockDatabaseService) SaveData(data string) error {
	args := m.Called(data)
	return args.Error(0)
}

// Benchmark functions for each component with different mocking strategies

func BenchmarkDataProcessor_Real(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dp := &RealDataProcessor{}
		dp.ProcessData("some-data")
	}
}

func BenchmarkDataProcessor_MockStub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dp := &MockDataProcessor{}
		dp.On("ProcessData", "some-data").Return("processed-data")
		dp.ProcessData("some-data")
	}
}

func BenchmarkDataProcessor_MockExpect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dp := &MockDataProcessor{}
		dp.On("ProcessData", "some-data").Return("processed-data")
		dp.ProcessData("some-data")
		dp.AssertExpectations(b)
	}
}

func BenchmarkDatabaseService_Real(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := &RealDatabaseService{}
		ds.SaveData("some-data")
	}
}

func BenchmarkDatabaseService_MockStub(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := &MockDatabaseService{}
		ds.On("SaveData", "some-data").Return(nil)
		ds.SaveData("some-data")
	}
}

func BenchmarkDatabaseService_MockExpect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ds := &MockDatabaseService{}
		ds.On("SaveData", "some-data").Return(nil)
		ds.SaveData("some-data")
		ds.AssertExpectations(b)
	}
}
