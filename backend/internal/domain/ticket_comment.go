package domain

import (
	"time"

	"github.com/google/uuid"
)

type TicketComment struct {
	ID              uuid.UUID `gorm:"column:comment_id;primaryKey" json:"comment_id"`
	TicketID        uuid.UUID `gorm:"column:ticket_id" json:"ticket_id"`
	TicketCreatedAt time.Time `gorm:"column:ticket_created_at" json:"ticket_created_at"`
	AuthorID        uuid.UUID `gorm:"column:author_id" json:"author_id"`
	Content         string    `gorm:"column:content" json:"content"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}
