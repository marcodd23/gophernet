package repository_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/marcodd23/gopernet/internal/models"
	"github.com/marcodd23/gopernet/internal/repository"
	"github.com/stretchr/testify/assert"
)

func setupTestRepo(t *testing.T) (*repository.MemoryRepository, func()) {
	// Create temporary state and report files
	stateFile, err := os.CreateTemp("", "state*.json")
	assert.NoError(t, err)
	reportFile, err := os.CreateTemp("", "report*.txt")
	assert.NoError(t, err)

	repo := repository.NewMemoryRepository(stateFile.Name(), reportFile.Name())

	// Return cleanup function to remove temp files
	return repo, func() {
		os.Remove(stateFile.Name())
		os.Remove(reportFile.Name())
	}
}

func TestMemoryRepository_LoadState(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Prepare initial burrows
	burrows := []*models.Burrow{
		{Name: "Burrow1", Depth: 1.5, Width: 1.0, Occupied: false, Age: 100},
		{Name: "Burrow2", Depth: 2.0, Width: 1.2, Occupied: true, Age: 50},
	}
	data, err := json.Marshal(burrows)
	assert.NoError(t, err)

	// Write initial state to file
	err = os.WriteFile(repo.GetStateFile(), data, 0644)
	assert.NoError(t, err)

	// Load state from file
	err = repo.LoadState()
	assert.NoError(t, err)

	// Check if burrows are loaded correctly using public API
	loadedBurrows := repo.GetAllBurrows()
	assert.Equal(t, 2, len(loadedBurrows))
	assert.Equal(t, "Burrow1", loadedBurrows[0].Name)
	assert.Equal(t, "Burrow2", loadedBurrows[1].Name)
}

func TestMemoryRepository_SaveState(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Add burrows to the repo using public API
	burrows := []*models.Burrow{
		{Name: "Burrow1", Depth: 1.5, Width: 1.0, Occupied: false, Age: 100},
		{Name: "Burrow2", Depth: 2.0, Width: 1.2, Occupied: true, Age: 50},
	}
	for _, b := range burrows {
		repo.AddBurrow(b)
	}

	// Save state to file
	err := repo.SaveState()
	assert.NoError(t, err)

	// Read the state back from the file
	data, err := os.ReadFile(repo.GetStateFile())
	assert.NoError(t, err)

	// Unmarshal the data and check the burrows
	var loadedBurrows []*models.Burrow
	err = json.Unmarshal(data, &loadedBurrows)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(loadedBurrows))
	assert.Equal(t, "Burrow1", loadedBurrows[0].Name)
	assert.Equal(t, "Burrow2", loadedBurrows[1].Name)
}

func TestMemoryRepository_RentBurrow(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Add a burrow to the repo using public API
	burrow := &models.Burrow{Name: "Burrow1", Depth: 1.5, Width: 1.0, Occupied: false, Age: 100}
	repo.AddBurrow(burrow)

	// Rent the burrow
	err := repo.RentBurrow("Burrow1")
	assert.NoError(t, err)

	// Check if the burrow is now occupied using public API
	loadedBurrows := repo.GetAllBurrows()
	assert.True(t, loadedBurrows[0].Occupied)

	// Try renting an already occupied burrow
	err = repo.RentBurrow("Burrow1")
	assert.Error(t, err)
}

func TestMemoryRepository_UpdateAllBurrows(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Add burrows to the repo using public API
	burrows := []*models.Burrow{
		{Name: "Burrow1", Depth: 1.0, Width: 1.0, Occupied: true, Age: 10},
		{Name: "Burrow2", Depth: 2.0, Width: 1.2, Occupied: false, Age: 20},
	}
	for _, b := range burrows {
		repo.AddBurrow(b)
	}

	// Update all burrows
	repo.UpdateAllBurrows()

	// Check that the depth of the occupied burrow increased
	loadedBurrows := repo.GetAllBurrows()
	assert.Greater(t, loadedBurrows[0].Depth, 1.0)
	// Check that the age increased for both burrows
	assert.Equal(t, 11, loadedBurrows[0].Age)
	assert.Equal(t, 21, loadedBurrows[1].Age)
}

func TestMemoryRepository_SaveReport(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	// Generate a report
	report := "Test Report Content"
	err := repo.SaveReport(report)
	assert.NoError(t, err)

	// Read the report back from the file
	data, err := os.ReadFile(repo.GetReportFile())
	assert.NoError(t, err)
	assert.Equal(t, report, string(data))
}
