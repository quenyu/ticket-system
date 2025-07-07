package repository

import (
	"ticket-system/backend/internal/domain"

	"gorm.io/gorm"
)

type TicketPriorityRepository struct {
	DB *gorm.DB
}

func NewTicketPriorityRepository(db *gorm.DB) *TicketPriorityRepository {
	return &TicketPriorityRepository{DB: db}
}

func (r *TicketPriorityRepository) List() ([]*domain.TicketPriority, error) {
	var priorities []*domain.TicketPriority
	if err := r.DB.Find(&priorities).Error; err != nil {
		return nil, err
	}
	return priorities, nil
}
