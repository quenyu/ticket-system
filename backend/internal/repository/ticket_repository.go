package repository

import (
	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

// Реализация методов интерфейса usecase.TicketRepository
// func (r *TicketRepository) Create(ticket *domain.Ticket) error { ... }
// ... и т.д.
