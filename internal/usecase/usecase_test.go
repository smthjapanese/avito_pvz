package usecase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/smthjapanese/avito_pvz/internal/pkg/jwt"
	"github.com/smthjapanese/avito_pvz/internal/repository"
	"github.com/smthjapanese/avito_pvz/internal/repository/mock"
)

func TestNewUseCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	pvzRepo := mock.NewMockPVZRepository(ctrl)
	receptionRepo := mock.NewMockReceptionRepository(ctrl)
	productRepo := mock.NewMockProductRepository(ctrl)

	repos := &repository.Repositories{
		User:      userRepo,
		PVZ:       pvzRepo,
		Reception: receptionRepo,
		Product:   productRepo,
	}

	tokenManager := jwt.NewManager("test-secret", time.Hour)

	useCases := NewUseCases(repos, tokenManager)
	assert.NotNil(t, useCases.User)
	assert.NotNil(t, useCases.PVZ)
	assert.NotNil(t, useCases.Reception)
	assert.NotNil(t, useCases.Product)
	
	_, ok := useCases.User.(*UserUseCase)
	assert.True(t, ok)

	_, ok = useCases.PVZ.(*PVZUseCase)
	assert.True(t, ok)

	_, ok = useCases.Reception.(*ReceptionUseCase)
	assert.True(t, ok)

	_, ok = useCases.Product.(*ProductUseCase)
	assert.True(t, ok)
}
