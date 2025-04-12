package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/repository"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
)

type PVZUseCase struct {
	pvzRepo       repository.PVZRepository
	receptionRepo repository.ReceptionRepository
	productRepo   repository.ProductRepository
}

func NewPVZUseCase(
	pvzRepo repository.PVZRepository,
	receptionRepo repository.ReceptionRepository,
	productRepo repository.ProductRepository,
) usecase.PVZUseCase {
	return &PVZUseCase{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
	}
}

func (uc *PVZUseCase) Create(ctx context.Context, city models.City) (*models.PVZ, error) {
	if !models.IsValidCity(city) {
		return nil, errors.ErrInvalidCity
	}

	pvz := models.NewPVZ(city)

	if err := uc.pvzRepo.Create(ctx, pvz); err != nil {
		return nil, err
	}

	return pvz, nil
}

func (uc *PVZUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.PVZ, error) {
	return uc.pvzRepo.GetByID(ctx, id)
}

func (uc *PVZUseCase) List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*usecase.PVZWithReceptions, error) {
	pvzs, err := uc.pvzRepo.List(ctx, startDate, endDate, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*usecase.PVZWithReceptions, 0, len(pvzs))
	for _, pvz := range pvzs {
		pvzWithReceptions, err := uc.getPVZWithReceptions(ctx, pvz)
		if err != nil {
			return nil, err
		}
		result = append(result, pvzWithReceptions)
	}

	return result, nil
}

func (uc *PVZUseCase) GetAll(ctx context.Context) ([]*models.PVZ, error) {
	return uc.pvzRepo.GetAll(ctx)
}

func (uc *PVZUseCase) getPVZWithReceptions(ctx context.Context, pvz *models.PVZ) (*usecase.PVZWithReceptions, error) {
	receptions, err := uc.receptionRepo.ListByPVZID(ctx, pvz.ID)
	if err != nil {
		return nil, err
	}
	
	receptionsWithProducts := make([]*usecase.ReceptionWithProducts, 0, len(receptions))
	for _, reception := range receptions {
		products, err := uc.productRepo.ListByReceptionID(ctx, reception.ID)
		if err != nil {
			return nil, err
		}
		receptionsWithProducts = append(receptionsWithProducts, &usecase.ReceptionWithProducts{
			Reception: reception,
			Products:  products,
		})
	}

	return &usecase.PVZWithReceptions{
		PVZ:        pvz,
		Receptions: receptionsWithProducts,
	}, nil
}
