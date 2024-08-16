package repository

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"sync"

	"github.com/marcodd23/gopernet/internal/models"
)

type MemoryRepository struct {
	burrows     map[string]*models.Burrow
	burrowsList []*models.Burrow // For preserving order
	mu          sync.RWMutex
	stateFile   string
	reportFile  string
}

func NewMemoryRepository(stateFile, reportFile string) *MemoryRepository {
	return &MemoryRepository{
		burrows:     make(map[string]*models.Burrow),
		burrowsList: make([]*models.Burrow, 0),
		stateFile:   stateFile,
		reportFile:  reportFile,
	}
}

func (s *MemoryRepository) GetAllBurrows() []*models.Burrow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	burrowsListCopy := make([]*models.Burrow, 0, len(s.burrowsList))
	for _, burrow := range s.burrowsList {
		// Create a deep copy of the burrow before returning
		// to avoid that the user of the repository could modify the
		// data in the storage
		copiedBurrow := &models.Burrow{
			Name:     burrow.Name,
			Depth:    burrow.Depth,
			Width:    burrow.Width,
			Occupied: burrow.Occupied,
			Age:      burrow.Age,
		}

		burrowsListCopy = append(burrowsListCopy, copiedBurrow)
	}

	return burrowsListCopy
}

func (s *MemoryRepository) RentBurrow(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	burrow, exists := s.burrows[name]
	if !exists {
		return errors.New("burrow not found")
	}

	if burrow.Occupied || burrow.HasCollapsed() {
		return errors.New("burrow not available")
	}

	burrow.Occupied = true

	return nil
}

func (s *MemoryRepository) UpdateAllBurrows() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, burrow := range s.burrowsList {
		burrow.UpdateDepth()
	}
}

func (s *MemoryRepository) LoadState() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.stateFile)
	if err != nil {
		return errors.WithMessage(err, "failed to read state file")
	}

	var burrows []*models.Burrow
	if err := json.Unmarshal(data, &burrows); err != nil {
		return errors.WithMessage(err, "failed to unmarshal state data")
	}

	// Clear existing data
	s.burrows = make(map[string]*models.Burrow)
	s.burrowsList = make([]*models.Burrow, 0)

	for _, burrow := range burrows {
		s.burrowsList = append(s.burrowsList, burrow)
		s.burrows[burrow.Name] = burrow
	}

	return nil
}

func (s *MemoryRepository) SaveState() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.MarshalIndent(s.burrowsList, "", "  ")
	if err != nil {
		return errors.WithMessage(err, "failed to marshal state")
	}

	if err := os.WriteFile(s.stateFile, data, 0644); err != nil {
		return errors.WithMessage(err, "failed to save state to file")
	}

	return nil
}

func (s *MemoryRepository) SaveReport(report string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.WriteFile(s.reportFile, []byte(report), 0644); err != nil {
		return errors.WithMessage(err, "failed to save report to file")
	}

	return nil
}

func (s *MemoryRepository) GetStateFile() string {
	return s.stateFile
}

func (s *MemoryRepository) GetReportFile() string {
	return s.reportFile
}

func (s *MemoryRepository) AddBurrow(burrow *models.Burrow) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.burrows[burrow.Name] = burrow
}
