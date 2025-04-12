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

func TestReceptionUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)

	uc := NewReceptionUseCase(pvzRepo, receptionRepo)

	pvzID := uuid.New()
	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(nil, errors.ErrOpenReceptionNotFound)

	receptionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, reception *models.Reception) error {
		assert.Equal(t, pvzID, reception.PVZID)
		assert.Equal(t, models.ReceptionStatusInProgress, reception.Status)
		return nil
	})

	reception, err := uc.Create(context.Background(), pvzID)
	require.NoError(t, err)
	assert.Equal(t, pvzID, reception.PVZID)
	assert.Equal(t, models.ReceptionStatusInProgress, reception.Status)
}

func TestReceptionUseCase_Create_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)

	uc := NewReceptionUseCase(pvzRepo, receptionRepo)

	pvzID := uuid.New()

	// ПВЗ не найден
	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(nil, errors.ErrPVZNotFound)

	_, err := uc.Create(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrPVZNotFound)
}

func TestReceptionUseCase_Create_OpenReceptionExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)

	uc := NewReceptionUseCase(pvzRepo, receptionRepo)

	pvzID := uuid.New()
	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	existingReception := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	// Уже есть открытая приемка
	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(existingReception, nil)

	_, err := uc.Create(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrOpenReceptionExists)
}

func TestReceptionUseCase_CloseLastReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)

	uc := NewReceptionUseCase(pvzRepo, receptionRepo)

	pvzID := uuid.New()
	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	reception := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	// Получение последней открытой приемки
	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(reception, nil)

	receptionRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, updatedReception *models.Reception) error {
		assert.Equal(t, reception.ID, updatedReception.ID)
		assert.Equal(t, models.ReceptionStatusClose, updatedReception.Status)
		return nil
	})

	result, err := uc.CloseLastReception(context.Background(), pvzID)
	require.NoError(t, err)
	assert.Equal(t, reception.ID, result.ID)
	assert.Equal(t, models.ReceptionStatusClose, result.Status)
}

func TestReceptionUseCase_CloseLastReception_PVZNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)

	uc := NewReceptionUseCase(pvzRepo, receptionRepo)

	pvzID := uuid.New()

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(nil, errors.ErrPVZNotFound)

	_, err := uc.CloseLastReception(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrPVZNotFound)
}

func TestReceptionUseCase_CloseLastReception_NoOpenReception(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)

	uc := NewReceptionUseCase(pvzRepo, receptionRepo)

	pvzID := uuid.New()
	pvz := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	pvzRepo.EXPECT().GetByID(gomock.Any(), pvzID).Return(pvz, nil)

	receptionRepo.EXPECT().GetLastOpenByPVZID(gomock.Any(), pvzID).Return(nil, errors.ErrOpenReceptionNotFound)

	_, err := uc.CloseLastReception(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrOpenReceptionNotFound)
}
