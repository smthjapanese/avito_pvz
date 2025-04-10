package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// PVZUseCase  интерфейс бизнес-логики с ПВЗ
type PVZUseCase interface {
	Create(ctx context.Context, city models.City) (*models.PVZ, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.PVZ, error)
	List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*PVZWithReceptions, error)
	GetAll(ctx context.Context) ([]*models.PVZ, error)
}

// PVZWithReceptions представляет ПВЗ с его приемками и товарами
type PVZWithReceptions struct {
	PVZ        *models.PVZ              `json:"pvz"`
	Receptions []*ReceptionWithProducts `json:"receptions"`
}

// ReceptionWithProducts представляет приемку с товарами
type ReceptionWithProducts struct {
	Reception *models.Reception `json:"reception"`
	Products  []*models.Product `json:"products"`
}
