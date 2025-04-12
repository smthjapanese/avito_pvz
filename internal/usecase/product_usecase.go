package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/repository"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
)

type ProductUseCase struct {
	pvzRepo       repository.PVZRepository
	receptionRepo repository.ReceptionRepository
	productRepo   repository.ProductRepository
}

func NewProductUseCase(
	pvzRepo repository.PVZRepository,
	receptionRepo repository.ReceptionRepository,
	productRepo repository.ProductRepository,
) usecase.ProductUseCase {
	return &ProductUseCase{
		pvzRepo:       pvzRepo,
		receptionRepo: receptionRepo,
		productRepo:   productRepo,
	}
}

func (uc *ProductUseCase) Create(ctx context.Context, productType models.ProductType, pvzID uuid.UUID) (*models.Product, error) {
	if !models.IsValidProductType(productType) {
		return nil, errors.ErrInvalidProductType
	}

	_, err := uc.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	reception, err := uc.receptionRepo.GetLastOpenByPVZID(ctx, pvzID)
	if err != nil {
		return nil, err
	}

	product := models.NewProduct(productType, reception.ID)

	if err := uc.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *ProductUseCase) DeleteLastFromReception(ctx context.Context, pvzID uuid.UUID) error {
	_, err := uc.pvzRepo.GetByID(ctx, pvzID)
	if err != nil {
		return err
	}

	reception, err := uc.receptionRepo.GetLastOpenByPVZID(ctx, pvzID)
	if err != nil {
		return err
	}

	product, err := uc.productRepo.GetLastByReceptionID(ctx, reception.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.ErrNoProductsToDelete
		}
		return err
	}
	
	if err := uc.productRepo.Delete(ctx, product.ID); err != nil {
		return err
	}

	return nil
}
