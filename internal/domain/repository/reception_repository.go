package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// ReceptionRepository представляет интерфейс для работы с хранилищем приемок
type ReceptionRepository interface {
	Create(ctx context.Context, reception *models.Reception) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Reception, error)
	GetLastByPVZID(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)
	GetLastOpenByPVZID(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)
	Update(ctx context.Context, reception *models.Reception) error
	ListByPVZID(ctx context.Context, pvzID uuid.UUID) ([]*models.Reception, error)
}
