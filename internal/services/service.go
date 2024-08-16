package services

import (
	"fmt"
	"math"

	"github.com/marcodd23/gopernet/internal/models"
	"github.com/marcodd23/gopernet/internal/repository"
)

type GopherService interface {
	LoadInitialState() error
	GetAllBurrows() []*models.Burrow
	RentBurrow(name string) error
	GenerateReport() (string, error)
	SaveState() error
	SaveReport() error
	UpdateBurrows()
}

type DefaultBurrowService struct {
	repo repository.StatefulRepository
}

func NewGopherNetService(repo repository.StatefulRepository) *DefaultBurrowService {
	return &DefaultBurrowService{repo: repo}
}

// LoadInitialState loads the initial state through the repository.
func (s *DefaultBurrowService) LoadInitialState() error {
	return s.repo.LoadState()
}

// GetAllBurrows returns all burrows through the repository.
func (s *DefaultBurrowService) GetAllBurrows() []*models.Burrow {
	return s.repo.GetAllBurrows()
}

// RentBurrow rents a burrow through the repository.
func (s *DefaultBurrowService) RentBurrow(name string) error {
	return s.repo.RentBurrow(name)
}

// GenerateReport generates a report of the current state of the burrows.
func (s *DefaultBurrowService) GenerateReport() (string, error) {
	burrows := s.repo.GetAllBurrows()

	totalDepth := 0.0
	availableBurrows := 0
	var largestBurrow, smallestBurrow *models.Burrow
	largestVolume := 0.0
	smallestVolume := math.MaxFloat64

	// Iterate over all burrows to calculate the report metrics.
	for _, burrow := range burrows {
		// Skip burrows with invalid dimensions
		if burrow.Width == 0 || burrow.Depth == 0 {
			continue
		}

		totalDepth += burrow.Depth

		if !burrow.Occupied && !burrow.HasCollapsed() {
			availableBurrows++
		}

		// Calculate the volume of the burrow (cylindrical volume formula: V = pi * r^2 * h).
		radius := burrow.Width / 2
		volume := math.Pi * radius * radius * burrow.Depth

		if volume > largestVolume {
			largestVolume = volume
			largestBurrow = burrow
		}
		if volume < smallestVolume {
			smallestVolume = volume
			smallestBurrow = burrow
		}
	}

	// Create the report string
	report := fmt.Sprintf("GopherNet Burrow Report\n\n")
	report += fmt.Sprintf("Total Depth of all Burrows: %.2f meters\n", totalDepth)
	report += fmt.Sprintf("Number of Available Burrows: %d\n", availableBurrows)

	if largestBurrow != nil {
		report += fmt.Sprintf("Largest Burrow by Volume: %s (%.2f cubic meters)\n", largestBurrow.Name, largestVolume)
	} else {
		report += "Largest Burrow by Volume: N/A\n"
	}

	if smallestBurrow != nil {
		report += fmt.Sprintf("Smallest Burrow by Volume: %s (%.2f cubic meters)\n", smallestBurrow.Name, smallestVolume)
	} else {
		report += "Smallest Burrow by Volume: N/A\n"
	}

	return report, nil
}

// SaveState instructs the repository to save the current state.
func (s *DefaultBurrowService) SaveState() error {
	return s.repo.SaveState()
}

// SaveReport generates the report and instructs the repository to save it.
func (s *DefaultBurrowService) SaveReport() error {
	report, err := s.GenerateReport()
	if err != nil {
		return err
	}

	return s.repo.SaveReport(report)
}

// UpdateBurrows triggers an update of all burrows through the repository.
func (s *DefaultBurrowService) UpdateBurrows() {
	s.repo.UpdateAllBurrows()
}
