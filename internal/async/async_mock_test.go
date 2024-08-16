package async_test

import (
	"github.com/marcodd23/gopernet/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockGopherService is a mock implementation of the GopherService interface.
type MockGopherService struct {
	mock.Mock
}

func (m *MockGopherService) LoadInitialState() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockGopherService) GetAllBurrows() []*models.Burrow {
	args := m.Called()
	return args.Get(0).([]*models.Burrow)
}

func (m *MockGopherService) RentBurrow(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockGopherService) GenerateReport() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockGopherService) SaveState() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockGopherService) SaveReport() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockGopherService) UpdateBurrows() {
	m.Called()
}
