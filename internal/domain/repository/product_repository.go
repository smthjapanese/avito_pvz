package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// ProductRepository представляет интерфейс для работы с хранилищем товаров
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	ListByReceptionID(ctx context.Context, receptionID uuid.UUID) ([]*models.Product, error)
	GetLastByReceptionID(ctx context.Context, receptionID uuid.UUID) (*models.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
