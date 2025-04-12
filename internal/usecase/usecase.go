package usecase

import (
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/jwt"
	repoProvider "github.com/smthjapanese/avito_pvz/internal/repository"
)

type UseCases struct {
	User      usecase.UserUseCase
	PVZ       usecase.PVZUseCase
	Reception usecase.ReceptionUseCase
	Product   usecase.ProductUseCase
}

func NewUseCases(repos *repoProvider.Repositories, tokenManager *jwt.Manager) *UseCases {
	return &UseCases{
		User:      NewUserUseCase(repos.User, tokenManager),
		PVZ:       NewPVZUseCase(repos.PVZ, repos.Reception, repos.Product),
		Reception: NewReceptionUseCase(repos.PVZ, repos.Reception),
		Product:   NewProductUseCase(repos.PVZ, repos.Reception, repos.Product),
	}
}
