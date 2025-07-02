package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID           uuid.UUID
	Title        string
	Description  string
	StatusID     int16
	PriorityID   int16
	CreatorID    uuid.UUID
	AssigneeID   uuid.UUID
	DepartmentID int16
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
