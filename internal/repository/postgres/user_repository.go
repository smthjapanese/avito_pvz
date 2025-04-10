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

type UserRepository struct {
	db *database.Database
	sb squirrel.StatementBuilderType
}

func NewUserRepository(db *database.Database) repository.UserRepository {
	return &UserRepository{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := r.sb.Insert("users").
		Columns("id", "email", "password_hash", "role").
		Values(user.ID, user.Email, user.PasswordHash, user.Role)

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

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := r.sb.Select("id", "email", "password_hash", "role", "created_at").
		From("users").
		Where(squirrel.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var user models.User
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRows(err) {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(errors.ErrDBQuery, fmt.Sprintf("failed to get user by ID: %v", err))
	}

	return &user, nil
}

// GetByEmail получает пользователя по email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := r.sb.Select("id", "email", "password_hash", "role", "created_at").
		From("users").
		Where(squirrel.Eq{"email": email})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var user models.User
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.IsNoRows(err) {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.Wrap(errors.ErrDBQuery, fmt.Sprintf("failed to get user by email: %v", err))
	}

	return &user, nil
}
