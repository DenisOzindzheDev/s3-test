package db

import "s3-test/internal/models"

type DB interface {
	CreateNewRentSession(rent models.Rent) error
	StartRentSession(int) error
	CompleteRentSession(int) error
	GetRentSessions(models.User) ([]models.Rent, error)
	GetRentStatus(int) (models.Rent, error)
	AddImageRecord(string, int) error
}
