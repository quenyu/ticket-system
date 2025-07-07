package usecase

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
)

type TicketFilter struct {
	StatusID     *int16
	AssigneeID   *string
	DepartmentID *int16
	Q            string
	Limit        int
	Offset       int
}

type TicketRepository interface {
	Create(ticket *domain.Ticket) error
	GetByID(id uuid.UUID) (*domain.Ticket, error)
	Update(ticket *domain.Ticket) error
	Delete(id uuid.UUID) error
	Search(filter TicketFilter) ([]*domain.Ticket, error)
}
