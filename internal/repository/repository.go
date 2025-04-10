package repository

import (
	"github.com/smthjapanese/avito_pvz/internal/domain/repository"
	"github.com/smthjapanese/avito_pvz/internal/pkg/database"
	"github.com/smthjapanese/avito_pvz/internal/repository/postgres"
)

type Repositories struct {
	User      repository.UserRepository
	PVZ       repository.PVZRepository
	Reception repository.ReceptionRepository
	Product   repository.ProductRepository
}

func NewRepositories(db *database.Database) *Repositories {
	return &Repositories{
		User:      postgres.NewUserRepository(db),
		PVZ:       postgres.NewPVZRepository(db),
		Reception: postgres.NewReceptionRepository(db),
		Product:   postgres.NewProductRepository(db),
	}
}
