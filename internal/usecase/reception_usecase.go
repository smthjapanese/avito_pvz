package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/repository"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
)

type ReceptionUseCase struct {
	pvzRepo       repository.PVZRepository
	receptionRepo repository.ReceptionRepository
}

func NewReceptionUseCase(
	pvzRepo repository.PVZRepository,
	receptionRepo repository.ReceptionRepository,
) usecase.ReceptionUseCase {
	return &ReceptionUseCase{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
	}
}

func (uc *ReceptionUseCase) Create(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	_, err := uc.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	lastOpenReception, err := uc.receptionRepo.GetLastOpenByPVZID(ctx, pvzID)
	if err == nil && lastOpenReception != nil {
		return nil, errors.ErrOpenReceptionExists
	}
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	reception := models.NewReception(pvzID)

	if err := uc.receptionRepo.Create(ctx, reception); err != nil {
		return nil, err
	}

	return reception, nil
}

func (uc *ReceptionUseCase) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	_, err := uc.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	reception, err := uc.receptionRepo.GetLastOpenByPVZID(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	reception.Close()

	if err := uc.receptionRepo.Update(ctx, reception); err != nil {
		return nil, err
	}

	return reception, nil
}
