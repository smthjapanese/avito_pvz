package usecase

import (
	"context"

	"github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// UserUseCase  интерфейс для работы с пользователями
type UserUseCase interface {
	Register(ctx context.Context, email, password string, role models.UserRole) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	DummyLogin(ctx context.Context, role models.UserRole) (string, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}
