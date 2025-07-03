package domain

import (
	"time"

	"github.com/google/uuid"
)

type TicketComment struct {
	ID              uuid.UUID `gorm:"column:comment_id;primaryKey"`
	TicketID        uuid.UUID `gorm:"column:ticket_id"`
	TicketCreatedAt time.Time `gorm:"column:ticket_created_at"`
	AuthorID        uuid.UUID `gorm:"column:author_id"`
	Content         string    `gorm:"column:content"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}
