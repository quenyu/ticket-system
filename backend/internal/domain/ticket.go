package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID           uuid.UUID  `gorm:"column:ticket_id;primaryKey"`
	Title        string     `gorm:"column:title"`
	Description  string     `gorm:"column:description"`
	StatusID     int16      `gorm:"column:status_id"`
	PriorityID   int16      `gorm:"column:priority_id"`
	CreatorID    uuid.UUID  `gorm:"column:creator_id"`
	AssigneeID   uuid.UUID  `gorm:"column:assignee_id"`
	DepartmentID int16      `gorm:"column:department_id"`
	CreatedAt    time.Time  `gorm:"column:created_at;primaryKey"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}
