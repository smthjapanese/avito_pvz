package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	// EmployeeRole представляет роль сотрудника ПВЗ
	EmployeeRole UserRole = "employee"
	// ModeratorRole представляет роль модератора
	ModeratorRole UserRole = "moderator"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // не включаем в JSON
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewUser(email string, passwordHash string, role UserRole) *User {
	return &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    time.Now(),
	}
}
