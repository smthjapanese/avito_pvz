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

func TestReceptionRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReceptionRepository(&database.Database{DB: db})

	pvzID := uuid.New()
	reception := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO receptions").
		WithArgs(reception.ID, reception.DateTime, reception.PVZID, reception.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), reception)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestReceptionRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReceptionRepository(&database.Database{DB: db})

	receptionID := uuid.New()
	pvzID := uuid.New()
	expectedReception := &models.Reception{
		ID:        receptionID,
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status", "created_at"}).
		AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.PVZID, expectedReception.Status, expectedReception.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM receptions").
		WithArgs(receptionID).
		WillReturnRows(rows)

	reception, err := repo.GetByID(context.Background(), receptionID)
	require.NoError(t, err)
	assert.Equal(t, expectedReception.ID, reception.ID)
	assert.Equal(t, expectedReception.PVZID, reception.PVZID)
	assert.Equal(t, expectedReception.Status, reception.Status)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestReceptionRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReceptionRepository(&database.Database{DB: db})

	receptionID := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM receptions").
		WithArgs(receptionID).
		WillReturnError(errors.ErrNoRows)

	_, err = repo.GetByID(context.Background(), receptionID)
	assert.ErrorIs(t, err, errors.ErrReceptionNotFound)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestReceptionRepository_GetLastByPVZID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReceptionRepository(&database.Database{DB: db})

	pvzID := uuid.New()
	expectedReception := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status", "created_at"}).
		AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.PVZID, expectedReception.Status, expectedReception.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM receptions").
		WithArgs(pvzID).
		WillReturnRows(rows)

	reception, err := repo.GetLastByPVZID(context.Background(), pvzID)
	require.NoError(t, err)
	assert.Equal(t, expectedReception.ID, reception.ID)
	assert.Equal(t, expectedReception.PVZID, reception.PVZID)
	assert.Equal(t, expectedReception.Status, reception.Status)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestReceptionRepository_GetLastOpenByPVZID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReceptionRepository(&database.Database{DB: db})

	pvzID := uuid.New()
	expectedReception := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now(),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status", "created_at"}).
		AddRow(expectedReception.ID, expectedReception.DateTime, expectedReception.PVZID, expectedReception.Status, expectedReception.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM receptions").
		WithArgs(pvzID, models.ReceptionStatusInProgress).
		WillReturnRows(rows)

	reception, err := repo.GetLastOpenByPVZID(context.Background(), pvzID)
	require.NoError(t, err)
	assert.Equal(t, expectedReception.ID, reception.ID)
	assert.Equal(t, expectedReception.PVZID, reception.PVZID)
	assert.Equal(t, expectedReception.Status, reception.Status)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestReceptionRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReceptionRepository(&database.Database{DB: db})

	reception := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now(),
		PVZID:     uuid.New(),
		Status:    models.ReceptionStatusClose,
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("UPDATE receptions").
		WithArgs(reception.Status, reception.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(context.Background(), reception)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestReceptionRepository_ListByPVZID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewReceptionRepository(&database.Database{DB: db})

	pvzID := uuid.New()
	reception1 := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now().Add(-12 * time.Hour),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusClose,
		CreatedAt: time.Now().Add(-12 * time.Hour),
	}

	reception2 := &models.Reception{
		ID:        uuid.New(),
		DateTime:  time.Now().Add(-6 * time.Hour),
		PVZID:     pvzID,
		Status:    models.ReceptionStatusInProgress,
		CreatedAt: time.Now().Add(-6 * time.Hour),
	}

	rows := sqlmock.NewRows([]string{"id", "date_time", "pvz_id", "status", "created_at"}).
		AddRow(reception1.ID, reception1.DateTime, reception1.PVZID, reception1.Status, reception1.CreatedAt).
		AddRow(reception2.ID, reception2.DateTime, reception2.PVZID, reception2.Status, reception2.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM receptions").
		WithArgs(pvzID).
		WillReturnRows(rows)

	receptions, err := repo.ListByPVZID(context.Background(), pvzID)
	require.NoError(t, err)
	assert.Len(t, receptions, 2)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
