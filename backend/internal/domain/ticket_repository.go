package domain

import "github.com/google/uuid"

// TicketRepository определяет методы для работы с тикетами
// (CRUD и поиск)
type TicketRepository interface {
	Create(ticket *Ticket) error
	GetByID(id uuid.UUID) (*Ticket, error)
	Update(ticket *Ticket) error
	Delete(id uuid.UUID) error
	Search(query string) ([]*Ticket, error)
} 