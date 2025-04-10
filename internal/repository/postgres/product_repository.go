package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/repository"
	"github.com/smthjapanese/avito_pvz/internal/pkg/database"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
)

type ProductRepository struct {
	db *database.Database
	sb squirrel.StatementBuilderType
}

func NewProductRepository(db *database.Database) repository.ProductRepository {
	return &ProductRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	query := r.sb.Insert("products").
		Columns("id", "date_time", "type", "reception_id").
		Values(product.ID, product.DateTime, product.Type, product.ReceptionID)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := r.sb.Select("id", "date_time", "type", "reception_id", "created_at").
		From("products").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var product models.Product
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
		&product.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRows(err) {
			return nil, errors.ErrProductNotFound
		}
		return nil, errors.Wrap(errors.ErrDBQuery, fmt.Sprintf("failed to get product by ID: %v", err))
	}

	return &product, nil
}

func (r *ProductRepository) ListByReceptionID(ctx context.Context, receptionID uuid.UUID) ([]*models.Product, error) {
	query := r.sb.Select("id", "date_time", "type", "reception_id", "created_at").
		From("products").
		Where(squirrel.Eq{"reception_id": receptionID}).
		OrderBy("date_time ASC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.DateTime,
			&product.Type,
			&product.ReceptionID,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return products, nil
}

func (r *ProductRepository) GetLastByReceptionID(ctx context.Context, receptionID uuid.UUID) (*models.Product, error) {
	query := r.sb.Select("id", "date_time", "type", "reception_id", "created_at").
		From("products").
		Where(squirrel.Eq{"reception_id": receptionID}).
		OrderBy("date_time DESC").
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var product models.Product
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&product.ID,
		&product.DateTime,
		&product.Type,
		&product.ReceptionID,
		&product.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRows(err) {
			return nil, errors.ErrProductNotFound
		}
		return nil, errors.Wrap(errors.ErrDBQuery, fmt.Sprintf("failed to get last product for reception: %v", err))
	}

	return &product, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := r.sb.Delete("products").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	result, err := r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
