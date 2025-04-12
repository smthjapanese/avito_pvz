package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
	"github.com/smthjapanese/avito_pvz/internal/repository/mock"
)

func TestProductUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewProductUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()
	receptionID := uuid.New()
	productType := models.ProductTypeElectronics

	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	reception := &models.Reception{
		ID:        receptionID,
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(reception, nil)

	productRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, product *models.Product) error {
		assert.Equal(t, receptionID, product.ReceptionID)
		assert.Equal(t, productType, product.Type)
		return nil
	})

	product, err := uc.Create(context.Background(), productType, pvzID)
	require.NoError(t, err)
	assert.Equal(t, receptionID, product.ReceptionID)
	assert.Equal(t, productType, product.Type)
}

func TestProductUseCase_Create_InvalidProductType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewProductUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()
	invalidProductType := models.ProductType("Invalid Type")

	_, err := uc.Create(context.Background(), invalidProductType, pvzID)
	assert.ErrorIs(t, err, errors.ErrInvalidProductType)
}

func TestProductUseCase_Create_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewProductUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()
	productType := models.ProductTypeElectronics

	// ПВЗ не найден
	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(nil, errors.ErrPVZNotFound)

	_, err := uc.Create(context.Background(), productType, pvzID)
	assert.ErrorIs(t, err, errors.ErrPVZNotFound)
}

func TestProductUseCase_Create_NoOpenReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewProductUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()
	productType := models.ProductTypeElectronics

	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(nil, errors.ErrOpenReceptionNotFound)

	_, err := uc.Create(context.Background(), productType, pvzID)
	assert.ErrorIs(t, err, errors.ErrOpenReceptionNotFound)
}

func TestProductUseCase_DeleteLastFromReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewProductUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()
	receptionID := uuid.New()
	productID := uuid.New()

	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	reception := &models.Reception{
		ID:        receptionID,
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	product := &models.Product{
		ID:          productID,
		DateTime:    time.Now(),
		Type:        models.ProductTypeElectronics,
		ReceptionID: receptionID,
		CreatedAt:   time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(reception, nil)

	productRepo.EXPECT().GetLastByReceptionID(gomock.Any(), receptionID).Return(product, nil)

	productRepo.EXPECT().Delete(gomock.Any(), productID).Return(nil)

	err := uc.DeleteLastFromReception(context.Background(), pvzID)
	require.NoError(t, err)
}

func TestProductUseCase_DeleteLastFromReception_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewProductUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(nil, errors.ErrPVZNotFound)

	err := uc.DeleteLastFromReception(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrPVZNotFound)
}

func TestProductUseCase_DeleteLastFromReception_NoOpenReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewProductUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()

	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(nil, errors.ErrOpenReceptionNotFound)

	err := uc.DeleteLastFromReception(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrOpenReceptionNotFound)
}

func TestProductUseCase_DeleteLastFromReception_NoProductsToDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewProductUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()
	receptionID := uuid.New()

	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	reception := &models.Reception{
		ID:        receptionID,
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(reception, nil)

	productRepo.EXPECT().GetLastByReceptionID(gomock.Any(), receptionID).Return(nil, errors.ErrProductNotFound)

	err := uc.DeleteLastFromReception(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrNoProductsToDelete)
}
