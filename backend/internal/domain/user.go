package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"column:user_id;primaryKey"`
	Username     string     `gorm:"column:username"`
	PasswordHash string     `gorm:"column:password_hash"`
	Email        string     `gorm:"column:email"`
	RoleID       uuid.UUID  `gorm:"column:role_id"`
	DepartmentID int16      `gorm:"column:department_id"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	LastLogin    *time.Time `gorm:"column:last_login"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}
