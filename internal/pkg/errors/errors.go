package errors

import (
	"errors"
	"fmt"
)

// Общие ошибки
var (
	ErrInternal      = errors.New("internal error")
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
)

// Ошибки для пользователей
var (
	ErrUserNotFound       = fmt.Errorf("user not found: %w", ErrNotFound)
	ErrUserAlreadyExists  = fmt.Errorf("user already exists: %w", ErrAlreadyExists)
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Ошибки для ПВЗ
var (
	ErrPVZNotFound = fmt.Errorf("pvz not found: %w", ErrNotFound)
	ErrInvalidCity = fmt.Errorf("invalid city: %w", ErrInvalidInput)
)

// Ошибки для приемок
var (
	ErrReceptionNotFound      = fmt.Errorf("reception not found: %w", ErrNotFound)
	ErrOpenReceptionNotFound  = fmt.Errorf("open reception not found: %w", ErrNotFound)
	ErrReceptionAlreadyClosed = errors.New("reception already closed")
	ErrOpenReceptionExists    = errors.New("open reception already exists")
)

// Ошибки для товаров
var (
	ErrProductNotFound    = fmt.Errorf("product not found: %w", ErrNotFound)
	ErrInvalidProductType = fmt.Errorf("invalid product type: %w", ErrInvalidInput)
	ErrNoProductsToDelete = errors.New("no products to delete")
)

// Ошибки базы данных
var (
	ErrDBConnection = errors.New("database connection error")
	ErrDBQuery      = errors.New("database query error")
	ErrNoRows       = errors.New("no rows in result set")
)

// IsNotFound проверяет, является ли ошибка типом "не найдено"
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsAlreadyExists проверяет, является ли ошибка типом "уже существует"
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}

// IsInvalidInput проверяет, является ли ошибка типом "неверный ввод"
func IsInvalidInput(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsUnauthorized проверяет, является ли ошибка типом "неавторизован"
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden проверяет, является ли ошибка типом "запрещено"
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsNoRows проверяет, является ли ошибка типом "нет строк"
func IsNoRows(err error) bool {
	return errors.Is(err, ErrNoRows) || (err != nil && err.Error() == "sql: no rows in result set")
}

// Wrap оборачивает ошибку с дополнительным сообщением
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}
