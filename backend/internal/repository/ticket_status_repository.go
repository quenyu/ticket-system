package repository

import (
	"ticket-system/backend/internal/domain"

	"gorm.io/gorm"
)

type TicketStatusRepository struct {
	DB *gorm.DB
}

func NewTicketStatusRepository(db *gorm.DB) *TicketStatusRepository {
	return &TicketStatusRepository{DB: db}
}

func (r *TicketStatusRepository) List() ([]*domain.TicketStatus, error) {
	var statuses []*domain.TicketStatus
	if err := r.DB.Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}
