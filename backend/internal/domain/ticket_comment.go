package domain

import (
	"time"

	"github.com/google/uuid"
)

type TicketComment struct {
	ID        uuid.UUID
	TicketID  uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	CreatedAt time.Time
}
