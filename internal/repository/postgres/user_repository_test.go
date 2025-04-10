package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/pkg/database"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(&database.Database{DB: db})

	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         models.EmployeeRole,
		CreatedAt:    time.Now(),
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Email, user.PasswordHash, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), user)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(&database.Database{DB: db})

	userID := uuid.New()
	expectedUser := &models.User{
		ID:           userID,
		Email:        "test@example.com",
		PasswordHash: "hash",
		Role:         models.EmployeeRole,
		CreatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "created_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PasswordHash, expectedUser.Role, expectedUser.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetByID(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.PasswordHash, user.PasswordHash)
	assert.Equal(t, expectedUser.Role, user.Role)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(&database.Database{DB: db})

	userID := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(userID).
		WillReturnError(errors.ErrNoRows)

	_, err = repo.GetByID(context.Background(), userID)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(&database.Database{DB: db})

	email := "test@example.com"
	expectedUser := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: "hash",
		Role:         models.EmployeeRole,
		CreatedAt:    time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "created_at"}).
		AddRow(expectedUser.ID, expectedUser.Email, expectedUser.PasswordHash, expectedUser.Role, expectedUser.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetByEmail(context.Background(), email)
	require.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.PasswordHash, user.PasswordHash)
	assert.Equal(t, expectedUser.Role, user.Role)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(&database.Database{DB: db})

	email := "nonexistent@example.com"

	mock.ExpectQuery("SELECT (.+) FROM users").
		WithArgs(email).
		WillReturnError(errors.ErrNoRows)

	_, err = repo.GetByEmail(context.Background(), email)
	assert.ErrorIs(t, err, errors.ErrUserNotFound)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
