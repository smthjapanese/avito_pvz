package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/smthjapanese/avito_pvz/internal/domain/models"
	"github.com/smthjapanese/avito_pvz/internal/domain/repository"
	"github.com/smthjapanese/avito_pvz/internal/pkg/database"
	"github.com/smthjapanese/avito_pvz/internal/pkg/errors"
)

type PVZRepository struct {
	db *database.Database
	sb squirrel.StatementBuilderType
}

func NewPVZRepository(db *database.Database) repository.PVZRepository {
	return &PVZRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PVZRepository) Create(ctx context.Context, pvz *models.PVZ) error {
	query := r.sb.Insert("pvzs").
		Columns("id", "registration_date", "city").
		Values(pvz.ID, pvz.RegistrationDate, pvz.City)

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

func (r *PVZRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PVZ, error) {
	query := r.sb.Select("id", "registration_date", "city", "created_at").
		From("pvzs").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var pvz models.PVZ
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&pvz.ID,
		&pvz.RegistrationDate,
		&pvz.City,
		&pvz.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRows(err) {
			return nil, errors.ErrPVZNotFound
		}
		return nil, errors.Wrap(errors.ErrDBQuery, fmt.Sprintf("failed to get PVZ by ID: %v", err))
	}

	return &pvz, nil
}

func (r *PVZRepository) List(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]*models.PVZ, error) {
	query := r.sb.Select("id", "registration_date", "city", "created_at").
		From("pvzs")

	if startDate != nil && endDate != nil {
		query = query.Where(squirrel.And{
			squirrel.GtOrEq{"registration_date": startDate},
			squirrel.LtOrEq{"registration_date": endDate},
		})
	} else if startDate != nil {
		query = query.Where(squirrel.GtOrEq{"registration_date": startDate})
	} else if endDate != nil {
		query = query.Where(squirrel.LtOrEq{"registration_date": endDate})
	}

	offset := (page - 1) * limit
	query = query.OrderBy("registration_date DESC").Limit(uint64(limit)).Offset(uint64(offset))

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var pvzs []*models.PVZ
	for rows.Next() {
		var pvz models.PVZ
		err := rows.Scan(
			&pvz.ID,
			&pvz.RegistrationDate,
			&pvz.City,
			&pvz.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		pvzs = append(pvzs, &pvz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return pvzs, nil
}

func (r *PVZRepository) GetAll(ctx context.Context) ([]*models.PVZ, error) {
	query := r.sb.Select("id", "registration_date", "city", "created_at").
		From("pvzs").
		OrderBy("registration_date DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var pvzs []*models.PVZ
	for rows.Next() {
		var pvz models.PVZ
		err := rows.Scan(
			&pvz.ID,
			&pvz.RegistrationDate,
			&pvz.City,
			&pvz.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		pvzs = append(pvzs, &pvz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return pvzs, nil
}
