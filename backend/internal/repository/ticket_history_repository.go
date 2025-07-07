package repository

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketHistoryRepository struct {
	DB *gorm.DB
}

func NewTicketHistoryRepository(db *gorm.DB) *TicketHistoryRepository {
	return &TicketHistoryRepository{DB: db}
}

func (r *TicketHistoryRepository) GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketHistory, error) {
	var history []*domain.TicketHistory
	err := r.DB.Where("ticket_id = ?", ticketID).Order("changed_at DESC").Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}
