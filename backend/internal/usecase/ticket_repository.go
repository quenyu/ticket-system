package usecase

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
)

type TicketRepository interface {
	Create(ticket *domain.Ticket) error
	GetByID(id uuid.UUID) (*domain.Ticket, error)
	Update(ticket *domain.Ticket) error
	Delete(id uuid.UUID) error
	Search(query string) ([]*domain.Ticket, error)
}
