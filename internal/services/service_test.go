package services_test

import (
	"github.com/marcodd23/gopernet/internal/models"
	"github.com/marcodd23/gopernet/internal/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockStatefulRepository struct {
	mock.Mock
}

func (m *MockStatefulRepository) LoadState() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStatefulRepository) GetAllBurrows() []*models.Burrow {
	args := m.Called()
	return args.Get(0).([]*models.Burrow)
}

func (m *MockStatefulRepository) RentBurrow(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockStatefulRepository) SaveState() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStatefulRepository) SaveReport(report string) error {
	args := m.Called(report)
	return args.Error(0)
}

func (m *MockStatefulRepository) UpdateAllBurrows() {
	m.Called()
}

func TestGopherNetService_LoadInitialState(t *testing.T) {
	mockRepo := new(MockStatefulRepository)
	service := services.NewGopherNetService(mockRepo)

	// Setup the mock expectation
	mockRepo.On("LoadState").Return(nil)

	// Call the method
	err := service.LoadInitialState()

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGopherNetService_GetAllBurrows(t *testing.T) {
	mockRepo := new(MockStatefulRepository)
	service := services.NewGopherNetService(mockRepo)

	// Setup the mock return value
	expectedBurrows := []*models.Burrow{
		{Name: "Burrow1", Depth: 1.5, Width: 1.0, Occupied: false, Age: 100},
	}
	mockRepo.On("GetAllBurrows").Return(expectedBurrows)

	// Call the method
	burrows := service.GetAllBurrows()

	// Assert the returned burrows
	assert.Equal(t, expectedBurrows, burrows)
	mockRepo.AssertExpectations(t)
}

func TestGopherNetService_RentBurrow(t *testing.T) {
	mockRepo := new(MockStatefulRepository)
	service := services.NewGopherNetService(mockRepo)

	// Setup the mock expectation
	mockRepo.On("RentBurrow", "Burrow1").Return(nil)

	// Call the method
	err := service.RentBurrow("Burrow1")

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGopherNetService_GenerateReport(t *testing.T) {
	mockRepo := new(MockStatefulRepository)
	service := services.NewGopherNetService(mockRepo)

	// Setup the mock return value
	burrows := []*models.Burrow{
		{Name: "Burrow1", Depth: 1.5, Width: 1.0, Occupied: false, Age: 100},
		{Name: "Burrow2", Depth: 2.0, Width: 1.2, Occupied: true, Age: 50},
	}
	mockRepo.On("GetAllBurrows").Return(burrows)

	// Call the method
	report, err := service.GenerateReport()

	// Assert the generated report
	assert.NoError(t, err)
	assert.Contains(t, report, "Total Depth of all Burrows: 3.50 meters")
	mockRepo.AssertExpectations(t)
}

func TestGopherNetService_SaveState(t *testing.T) {
	mockRepo := new(MockStatefulRepository)
	service := services.NewGopherNetService(mockRepo)

	// Setup the mock expectation
	mockRepo.On("SaveState").Return(nil)

	// Call the method
	err := service.SaveState()

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGopherNetService_SaveReport(t *testing.T) {
	mockRepo := new(MockStatefulRepository)
	service := services.NewGopherNetService(mockRepo)

	// Setup the mock return value and expectations
	burrows := []*models.Burrow{
		{Name: "Burrow1", Depth: 1.5, Width: 1.0, Occupied: false, Age: 100},
	}
	mockRepo.On("GetAllBurrows").Return(burrows)
	mockRepo.On("SaveReport", mock.AnythingOfType("string")).Return(nil)

	// Call the method
	err := service.SaveReport()

	// Assert expectations
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGopherNetService_UpdateBurrows(t *testing.T) {
	mockRepo := new(MockStatefulRepository)
	service := services.NewGopherNetService(mockRepo)

	// Setup the mock expectation
	mockRepo.On("UpdateAllBurrows").Return()

	// Call the method
	service.UpdateBurrows()

	// Assert expectations
	mockRepo.AssertExpectations(t)
}
