package service

// Storage represents an external storage interface
type Storage interface {
	GetData() ([]byte, error)
}

// Service encapsulates the service logic
type Service struct {
	storage Storage
}

// NewService creates a new Service instance
func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// Process Data Processing function
func (s *Service) Process() ([]byte, error) {
	data, err := s.storage.GetData()
	if err != nil {
		return nil, err
	}
	// Add more processing logic if needed
	return data, nil
}
