package repository

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"sync"

	"github.com/marcodd23/gopernet/internal/models"
)

type MemoryRepository struct {
	burrows    map[string]*models.Burrow
	mu         sync.RWMutex
	stateFile  string
	reportFile string
}

func NewMemoryRepository(stateFile, reportFile string) *MemoryRepository {
	return &MemoryRepository{
		burrows:    make(map[string]*models.Burrow),
		stateFile:  stateFile,
		reportFile: reportFile,
	}
}

func (s *MemoryRepository) GetAllBurrows() []*models.Burrow {
	s.mu.RLock()
	defer s.mu.RUnlock()

	burrows := make([]*models.Burrow, 0, len(s.burrows))
	for _, burrow := range s.burrows {
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

		burrows = append(burrows, copiedBurrow)
	}

	return burrows
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

	for _, burrow := range s.burrows {
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

	for _, burrow := range burrows {
		s.burrows[burrow.Name] = burrow
	}

	return nil
}

func (s *MemoryRepository) SaveState() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert the map to a slice of Burrow pointers
	burrowSlice := make([]*models.Burrow, 0, len(s.burrows))
	for _, burrow := range s.burrows {
		burrowSlice = append(burrowSlice, burrow)
	}

	data, err := json.MarshalIndent(burrowSlice, "", "  ")
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
