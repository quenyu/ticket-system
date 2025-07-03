package domain

import "github.com/google/uuid"

type Role struct {
	ID   uuid.UUID `gorm:"column:role_id;primaryKey"`
	Name string    `gorm:"column:name"`
}
