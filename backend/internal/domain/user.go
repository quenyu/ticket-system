package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"column:user_id;primaryKey" json:"user_id"`
	Username     string     `gorm:"column:username" json:"username"`
	PasswordHash string     `gorm:"column:password_hash" json:"-"`
	Email        string     `gorm:"column:email" json:"email"`
	RoleID       uuid.UUID  `gorm:"column:role_id" json:"role_id"`
	DepartmentID int16      `gorm:"column:department_id" json:"department_id"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	LastLogin    *time.Time `gorm:"column:last_login" json:"last_login,omitempty"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}
