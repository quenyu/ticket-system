package repository

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketHistoryRepository struct {
	db *gorm.DB
}

func NewTicketHistoryRepository(db *gorm.DB) *TicketHistoryRepository {
	return &TicketHistoryRepository{db: db}
}

func (r *TicketHistoryRepository) GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketHistory, error) {
	var history []*domain.TicketHistory
	err := r.db.Where("ticket_id = ?", ticketID).Order("changed_at DESC").Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (r *TicketHistoryRepository) Create(history *domain.TicketHistory) error {
	if history.ID == uuid.Nil {
		history.ID = uuid.New()
	}
	return r.db.Create(history).Error
}
