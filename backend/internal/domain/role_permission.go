package domain

import "github.com/google/uuid"

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"column:role_id;primaryKey"`
	PermissionID uuid.UUID `gorm:"column:permission_id;primaryKey"`
}
