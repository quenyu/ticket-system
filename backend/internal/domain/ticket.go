package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID           uuid.UUID  `gorm:"column:ticket_id;primaryKey" json:"ticket_id"`
	Title        string     `gorm:"column:title" json:"title"`
	Description  string     `gorm:"column:description" json:"description"`
	StatusID     int16      `gorm:"column:status_id" json:"status_id"`
	PriorityID   int16      `gorm:"column:priority_id" json:"priority_id"`
	CreatorID    uuid.UUID  `gorm:"column:creator_id" json:"creator_id"`
	AssigneeID   uuid.UUID  `gorm:"column:assignee_id" json:"assignee_id"`
	DepartmentID int16      `gorm:"column:department_id" json:"department_id"`
	CreatedAt    time.Time  `gorm:"column:created_at;primaryKey" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
