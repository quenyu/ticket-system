package usecase

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
)

type TicketCommentRepository interface {
	GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketComment, error)
	Create(comment *domain.TicketComment) error
}
