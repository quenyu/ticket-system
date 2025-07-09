package usecase

import (
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"

	"github.com/google/uuid"
)

type TicketRepository interface {
	Create(ticket *domain.Ticket) error
	GetByID(id uuid.UUID) (*domain.Ticket, error)
	Update(ticket *domain.Ticket) error
	Delete(id uuid.UUID) error
	Search(filter model.TicketFilter) ([]*domain.Ticket, error)
}
