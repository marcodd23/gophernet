package repository

import "github.com/marcodd23/gopernet/internal/models"

type Repository interface {
	GetAllBurrows() []*models.Burrow
	RentBurrow(name string) error
	UpdateAllBurrows()
}

type StatefulRepository interface {
	Repository
	LoadState() error
	SaveState() error
	SaveReport(report string) error
}
