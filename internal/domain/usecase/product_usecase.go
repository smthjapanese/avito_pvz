package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// ProductUseCase  интерфейс для работы с товарами
type ProductUseCase interface {
	Create(ctx context.Context, productType models.ProductType, pvzID uuid.UUID) (*models.Product, error)
	DeleteLastFromReception(ctx context.Context, pvzID uuid.UUID) error
}
