package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	Email        string
	RoleID       uuid.UUID
	DepartmentID int16
	CreatedAt    time.Time
	LastLogin    *time.Time
	DeletedAt    *time.Time
}
