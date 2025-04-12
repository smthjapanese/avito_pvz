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
	"github.com/smthjapanese/avito_pvz/internal/pkg/jwt"
	"github.com/smthjapanese/avito_pvz/internal/repository/mock"
)

func TestUserUseCase_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	tokenManager := jwt.NewManager("test-secret", time.Hour)

	uc := NewUserUseCase(userRepo, tokenManager)

	email := "test@example.com"
	password := "password"
	role := models.EmployeeRole

	userRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(nil, errors.ErrUserNotFound)

	userRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, user *models.User) error {
		assert.Equal(t, email, user.Email)
		assert.NotEmpty(t, user.PasswordHash)
		assert.Equal(t, role, user.Role)
		return nil
	})

	user, err := uc.Register(context.Background(), email, password, role)
	require.NoError(t, err)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, role, user.Role)
}

func TestUserUseCase_Register_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	tokenManager := jwt.NewManager("test-secret", time.Hour)

	uc := NewUserUseCase(userRepo, tokenManager)

	email := "test@example.com"
	password := "password"
	role := models.EmployeeRole

	existingUser := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: "hash",
		Role:         role,
		CreatedAt:    time.Now(),
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(existingUser, nil)

	_, err := uc.Register(context.Background(), email, password, role)
	assert.ErrorIs(t, err, errors.ErrUserAlreadyExists)
}

func TestUserUseCase_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	tokenManager := jwt.NewManager("test-secret", time.Hour)

	uc := NewUserUseCase(userRepo, tokenManager)

	email := "test@example.com"
	password := "password"

	hashedPassword, err := uc.(*UserUseCase).hashPassword(password)
	require.NoError(t, err)

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         models.EmployeeRole,
		CreatedAt:    time.Now(),
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(user, nil)

	token, err := uc.Login(context.Background(), email, password)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestUserUseCase_Login_InvalidCredentials(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	tokenManager := jwt.NewManager("test-secret", time.Hour)

	uc := NewUserUseCase(userRepo, tokenManager)

	email := "test@example.com"
	password := "password"
	wrongPassword := "wrong-password"

	hashedPassword, err := uc.(*UserUseCase).hashPassword(password)
	require.NoError(t, err)

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         models.EmployeeRole,
		CreatedAt:    time.Now(),
	}

	userRepo.EXPECT().GetByEmail(gomock.Any(), email).Return(user, nil)

	_, err = uc.Login(context.Background(), email, wrongPassword)
	assert.ErrorIs(t, err, errors.ErrInvalidCredentials)
}

func TestUserUseCase_DummyLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock.NewMockUserRepository(ctrl)
	tokenManager := jwt.NewManager("test-secret", time.Hour)

	uc := NewUserUseCase(userRepo, tokenManager)

	role := models.EmployeeRole

	token, err := uc.DummyLogin(context.Background(), role)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := tokenManager.ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, role, claims.Role)
	assert.Contains(t, claims.Email, "dummy_")
}
