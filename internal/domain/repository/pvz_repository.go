package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// PVZRepository представляет интерфейс для работы с хранилищем ПВЗ
type PVZRepository interface {
	Create(ctx context.Context, pvz *models.PVZ) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.PVZ, error)
	List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*models.PVZ, error)
	GetAll(ctx context.Context) ([]*models.PVZ, error)
}
