package usecase

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
)

type TicketHistoryRepository interface {
	GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketHistory, error)
}
