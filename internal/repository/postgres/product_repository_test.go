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

func TestProductRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(&database.Database{DB: db})

	receptionID := uuid.New()
	product := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        models.ProductTypeElectronics,
		ReceptionID: receptionID,
		CreatedAt:   time.Now(),
	}

	mock.ExpectExec("INSERT INTO products").
		WithArgs(product.ID, product.DateTime, product.Type, product.ReceptionID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), product)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestProductRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(&database.Database{DB: db})

	productID := uuid.New()
	receptionID := uuid.New()
	expectedProduct := &models.Product{
		ID:          productID,
		DateTime:    time.Now(),
		Type:        models.ProductTypeElectronics,
		ReceptionID: receptionID,
		CreatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id", "created_at"}).
		AddRow(expectedProduct.ID, expectedProduct.DateTime, expectedProduct.Type, expectedProduct.ReceptionID, expectedProduct.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM products").
		WithArgs(productID).
		WillReturnRows(rows)

	product, err := repo.GetByID(context.Background(), productID)
	require.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Type, product.Type)
	assert.Equal(t, expectedProduct.ReceptionID, product.ReceptionID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(&database.Database{DB: db})

	productID := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM products").
		WithArgs(productID).
		WillReturnError(errors.ErrNoRows)

	_, err = repo.GetByID(context.Background(), productID)
	assert.ErrorIs(t, err, errors.ErrProductNotFound)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestProductRepository_ListByReceptionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(&database.Database{DB: db})

	receptionID := uuid.New()
	product1 := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now().Add(-2 * time.Hour),
		Type:        models.ProductTypeElectronics,
		ReceptionID: receptionID,
		CreatedAt:   time.Now().Add(-2 * time.Hour),
	}

	product2 := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now().Add(-1 * time.Hour),
		Type:        models.ProductTypeClothes,
		ReceptionID: receptionID,
		CreatedAt:   time.Now().Add(-1 * time.Hour),
	}

	rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id", "created_at"}).
		AddRow(product1.ID, product1.DateTime, product1.Type, product1.ReceptionID, product1.CreatedAt).
		AddRow(product2.ID, product2.DateTime, product2.Type, product2.ReceptionID, product2.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM products").
		WithArgs(receptionID).
		WillReturnRows(rows)

	products, err := repo.ListByReceptionID(context.Background(), receptionID)
	require.NoError(t, err)
	assert.Len(t, products, 2)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestProductRepository_GetLastByReceptionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(&database.Database{DB: db})

	receptionID := uuid.New()
	expectedProduct := &models.Product{
		ID:          uuid.New(),
		DateTime:    time.Now(),
		Type:        models.ProductTypeElectronics,
		ReceptionID: receptionID,
		CreatedAt:   time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "date_time", "type", "reception_id", "created_at"}).
		AddRow(expectedProduct.ID, expectedProduct.DateTime, expectedProduct.Type, expectedProduct.ReceptionID, expectedProduct.CreatedAt)

	mock.ExpectQuery("SELECT (.+) FROM products").
		WithArgs(receptionID).
		WillReturnRows(rows)

	product, err := repo.GetLastByReceptionID(context.Background(), receptionID)
	require.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Type, product.Type)
	assert.Equal(t, expectedProduct.ReceptionID, product.ReceptionID)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestProductRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(&database.Database{DB: db})

	productID := uuid.New()

	mock.ExpectExec("DELETE FROM products").
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), productID)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestProductRepository_Delete_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewProductRepository(&database.Database{DB: db})

	productID := uuid.New()

	mock.ExpectExec("DELETE FROM products").
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.Delete(context.Background(), productID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
