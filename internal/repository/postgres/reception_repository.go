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

type ReceptionRepository struct {
	db *database.Database
	sb squirrel.StatementBuilderType
}

func NewReceptionRepository(db *database.Database) repository.ReceptionRepository {
	return &ReceptionRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ReceptionRepository) Create(ctx context.Context, reception *models.Reception) error {
	query := r.sb.Insert("receptions").
		Columns("id", "date_time", "pvz_id", "status").
		Values(reception.ID, reception.DateTime, reception.PVZID, reception.Status)

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

func (r *ReceptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Reception, error) {
	query := r.sb.Select("id", "date_time", "pvz_id", "status", "created_at").
		From("receptions").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var reception models.Reception
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
		&reception.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRows(err) {
			return nil, errors.ErrReceptionNotFound
		}
		return nil, errors.Wrap(errors.ErrDBQuery, fmt.Sprintf("failed to get reception by ID: %v", err))
	}

	return &reception, nil
}

func (r *ReceptionRepository) GetLastByPVZID(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	query := r.sb.Select("id", "date_time", "pvz_id", "status", "created_at").
		From("receptions").
		Where(squirrel.Eq{"pvz_id": pvzID}).
		OrderBy("date_time DESC").
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var reception models.Reception
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
		&reception.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRows(err) {
			return nil, errors.ErrReceptionNotFound
		}
		return nil, errors.Wrap(errors.ErrDBQuery, fmt.Sprintf("failed to get last reception for PVZ: %v", err))
	}

	return &reception, nil
}

func (r *ReceptionRepository) GetLastOpenByPVZID(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	query := r.sb.Select("id", "date_time", "pvz_id", "status", "created_at").
		From("receptions").
		Where(squirrel.And{
			squirrel.Eq{"pvz_id": pvzID},
			squirrel.Eq{"status": models.ReceptionStatusInProgress},
		}).
		OrderBy("date_time DESC").
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var reception models.Reception
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&reception.ID,
		&reception.DateTime,
		&reception.PVZID,
		&reception.Status,
		&reception.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRows(err) {
			return nil, errors.ErrOpenReceptionNotFound
		}
		return nil, errors.Wrap(errors.ErrDBQuery, fmt.Sprintf("failed to get last open reception for PVZ: %v", err))
	}

	return &reception, nil
}

func (r *ReceptionRepository) Update(ctx context.Context, reception *models.Reception) error {
	query := r.sb.Update("receptions").
		Set("status", reception.Status).
		Where(squirrel.Eq{"id": reception.ID})

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
		return fmt.Errorf("reception not found")
	}

	return nil
}

func (r *ReceptionRepository) ListByPVZID(ctx context.Context, pvzID uuid.UUID) ([]*models.Reception, error) {
	query := r.sb.Select("id", "date_time", "pvz_id", "status", "created_at").
		From("receptions").
		Where(squirrel.Eq{"pvz_id": pvzID}).
		OrderBy("date_time DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var receptions []*models.Reception
	for rows.Next() {
		var reception models.Reception
		err := rows.Scan(
			&reception.ID,
			&reception.DateTime,
			&reception.PVZID,
			&reception.Status,
			&reception.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		receptions = append(receptions, &reception)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return receptions, nil
}
