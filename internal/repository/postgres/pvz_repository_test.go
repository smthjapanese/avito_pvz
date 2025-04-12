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

func TestPVZRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPVZRepository(&database.Database{DB: db})

	pvz := &models.PVZ{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	mock.ExpectExec("INSERT INTO pvzs").
		WithArgs(pvz.ID, pvz.RegistrationDate, pvz.City).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), pvz)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPVZRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPVZRepository(&database.Database{DB: db})

	pvzID := uuid.New()
	expectedPVZ := &models.PVZ{
		ID:               pvzID,
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
		CreatedAt:        time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "registration_date", "city", "created_at"}).
		AddRow(expectedPVZ.ID, expectedPVZ.RegistrationDate, expectedPVZ.City, expectedPVZ.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM pvzs").
		WithArgs(pvzID).
		WillReturnRows(rows)

	pvz, err := repo.GetByID(context.Background(), pvzID)
	require.NoError(t, err)
	assert.Equal(t, expectedPVZ.ID, pvz.ID)
	assert.Equal(t, expectedPVZ.City, pvz.City)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPVZRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPVZRepository(&database.Database{DB: db})

	pvzID := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM pvzs").
		WithArgs(pvzID).
		WillReturnError(errors.ErrNoRows)

	_, err = repo.GetByID(context.Background(), pvzID)
	assert.ErrorIs(t, err, errors.ErrPVZNotFound)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPVZRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPVZRepository(&database.Database{DB: db})

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

	rows := sqlmock.NewRows([]string{"id", "registration_date", "city", "created_at"}).
		AddRow(pvz1.ID, pvz1.RegistrationDate, pvz1.City, pvz1.CreatedAt).
		AddRow(pvz2.ID, pvz2.RegistrationDate, pvz2.City, pvz2.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM pvzs").
		WillReturnRows(rows)

	pvzs, err := repo.List(context.Background(), &startDate, &endDate, page, limit)
	require.NoError(t, err)
	assert.Len(t, pvzs, 2)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPVZRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewPVZRepository(&database.Database{DB: db})

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

	rows := sqlmock.NewRows([]string{"id", "registration_date", "city", "created_at"}).
		AddRow(pvz1.ID, pvz1.RegistrationDate, pvz1.City, pvz1.CreatedAt).
		AddRow(pvz2.ID, pvz2.RegistrationDate, pvz2.City, pvz2.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM pvzs").
		WillReturnRows(rows)

	pvzs, err := repo.GetAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, pvzs, 2)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
