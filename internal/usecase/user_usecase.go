package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/repository"
	"github.com/smthjapanese/avito_pvz/internal/domain/usecase"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
	"github.com/smthjapanese/avito_pvz/internal/pkg/jwt"
	"github.com/smthjapanese/avito_pvz/internal/pkg/password"
)

type UserUseCase struct {
	userRepo     repository.UserRepository
	tokenManager *jwt.Manager
}

func NewUserUseCase(userRepo repository.UserRepository, tokenManager *jwt.Manager) usecase.UserUseCase {
	return &UserUseCase{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, email, plainPassword string, role models.UserRole) (*models.User, error) {
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}
	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	hashedPassword, err := password.Hash(plainPassword, password.DefaultParams())
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password")
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         role,
		CreatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) Login(ctx context.Context, email, plainPassword string) (string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFound(err) {
			return "", errors.ErrInvalidCredentials
		}
		return "", err
	}

	isValid, err := password.Verify(plainPassword, user.PasswordHash)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to verify password")
	}
	if !isValid {
		return "", errors.ErrInvalidCredentials
	}

	token, err := uc.tokenManager.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to generate token")
	}

	return token, nil
}

func (uc *UserUseCase) DummyLogin(ctx context.Context, role models.UserRole) (string, error) {
	if role != models.EmployeeRole && role != models.ModeratorRole {
		return "", errors.ErrInvalidInput
	}

	token, err := uc.tokenManager.GenerateDummyToken(role)
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to generate dummy token")
	}

	return token, nil
}

func (uc *UserUseCase) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	claims, err := uc.tokenManager.ParseToken(tokenString)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	if claims.Email[:6] == "dummy_" {
		return &models.User{
			ID:    userID,
			Email: claims.Email,
			Role:  claims.Role,
		}, nil
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUnauthorized
	}

	return user, nil
}

func (uc *UserUseCase) hashPassword(plainPassword string) (string, error) {
	return password.Hash(plainPassword, password.DefaultParams())
}
