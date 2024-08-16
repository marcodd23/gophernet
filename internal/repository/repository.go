package repository

import "github.com/marcodd23/gopernet/internal/models"

type Repository interface {
	GetAllBurrows() []*models.Burrow
	RentBurrow(name string) error
	UpdateAllBurrows()
	AddBurrow(burrow *models.Burrow)
}

type StatefulRepository interface {
	Repository
	LoadState() error
	SaveState() error
	SaveReport(report string) error
	GetStateFile() string
	GetReportFile() string
}
