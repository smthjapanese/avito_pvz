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

func TestPVZUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewPVZUseCase(pvzRepo, receptionRepo, productRepo)

	city := models.CityMoscow

	pvzRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, pvz *models.PVZ) error {
		assert.Equal(t, city, pvz.City)
		return nil
	})

	pvz, err := uc.Create(context.Background(), city)
	require.NoError(t, err)
	assert.Equal(t, city, pvz.City)
}

func TestPVZUseCase_Create_InvalidCity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewPVZUseCase(pvzRepo, receptionRepo, productRepo)

	invalidCity := models.City("Invalid City")

	_, err := uc.Create(context.Background(), invalidCity)
	assert.ErrorIs(t, err, errors.ErrInvalidCity)
}

func TestPVZUseCase_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewPVZUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()
	expectedPVZ := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(expectedPVZ, nil)

	pvz, err := uc.GetByID(context.Background(), pvzID)
	require.NoError(t, err)
	assert.Equal(t, expectedPVZ, pvz)
}

func TestPVZUseCase_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewPVZUseCase(pvzRepo, receptionRepo, productRepo)

	pvzID := uuid.New()

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(nil, errors.ErrPVZNotFound)

	_, err := uc.GetByID(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrPVZNotFound)
}

func TestPVZUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewPVZUseCase(pvzRepo, receptionRepo, productRepo)

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	page := 1
	limit := 10

	pvz1 := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now().Add(-12 * time.Hour),
		City:             models.CityMoscow,
		CreatedAt:        time.Now().Add(-12 * time.Hour),
	}

	pvz2 := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now().Add(-6 * time.Hour),
		City:             models.CitySaintPetersburg,
		CreatedAt:        time.Now().Add(-6 * time.Hour),
	}

	pvzs := []*models.PVZ{pvz1, pvz2}

	pvzRepo.EXPECT().List(gomock.Any(), &startDate, &endDate, page, limit).Return(pvzs, nil)

	for _, pvz := range pvzs {
		reception1 := &models.Reception{
			ID:        uuid.New(),
			DateTime:  time.Now().Add(-10 * time.Hour),
			PVZID:     pvz.ID,
			Status:    models.ReceptionStatusClose,
			CreatedAt: time.Now().Add(-10 * time.Hour),
		}

		reception2 := &models.Reception{
			ID:        uuid.New(),
			DateTime:  time.Now().Add(-5 * time.Hour),
			PVZID:     pvz.ID,
			Status:    models.ReceptionStatusInProgress,
			CreatedAt: time.Now().Add(-5 * time.Hour),
		}

		receptions := []*models.Reception{reception1, reception2}

		receptionRepo.EXPECT().ListByPVZID(gomock.Any(), pvz.ID).Return(receptions, nil)

		for _, reception := range receptions {
			product1 := &models.Product{
				ID:          uuid.New(),
				DateTime:    time.Now().Add(-9 * time.Hour),
				Type:        models.ProductTypeElectronics,
				ReceptionID: reception.ID,
				CreatedAt:   time.Now().Add(-9 * time.Hour),
			}

			product2 := &models.Product{
				ID:          uuid.New(),
				DateTime:    time.Now().Add(-8 * time.Hour),
				Type:        models.ProductTypeClothes,
				ReceptionID: reception.ID,
				CreatedAt:   time.Now().Add(-8 * time.Hour),
			}

			products := []*models.Product{product1, product2}

			productRepo.EXPECT().ListByReceptionID(gomock.Any(), reception.ID).Return(products, nil)
		}
	}

	result, err := uc.List(context.Background(), &startDate, &endDate, page, limit)
	require.NoError(t, err)
	assert.Len(t, result, 2)

	for i, pvzWithReceptions := range result {
		assert.Equal(t, pvzs[i], pvzWithReceptions.PVZ)
		assert.Len(t, pvzWithReceptions.Receptions, 2)

		for _, receptionWithProducts := range pvzWithReceptions.Receptions {
			assert.NotNil(t, receptionWithProducts.Reception)
			assert.Len(t, receptionWithProducts.Products, 2)
		}
	}
}

func TestPVZUseCase_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	uc := NewPVZUseCase(pvzRepo, receptionRepo, productRepo)

	pvz1 := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now().Add(-12 * time.Hour),
		City:             models.CityMoscow,
		CreatedAt:        time.Now().Add(-12 * time.Hour),
	}

	pvz2 := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now().Add(-6 * time.Hour),
		City:             models.CitySaintPetersburg,
		CreatedAt:        time.Now().Add(-6 * time.Hour),
	}

	pvzs := []*models.PVZ{pvz1, pvz2}

	pvzRepo.EXPECT().GetAll(gomock.Any()).Return(pvzs, nil)

	result, err := uc.GetAll(context.Background())
	require.NoError(t, err)
	assert.Equal(t, pvzs, result)
}
