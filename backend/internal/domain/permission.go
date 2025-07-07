package domain

import "github.com/google/uuid"

type Permission struct {
	ID   uuid.UUID `gorm:"column:permission_id;primaryKey"`
	Name string    `gorm:"column:name"`
}
