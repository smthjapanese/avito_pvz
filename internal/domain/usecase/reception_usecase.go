package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
)

// ReceptionUseCase интерфейс для работы с приемками
type ReceptionUseCase interface {
	Create(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)
}
