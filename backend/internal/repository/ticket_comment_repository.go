package repository

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketCommentRepository struct {
	DB *gorm.DB
}

func NewTicketCommentRepository(db *gorm.DB) *TicketCommentRepository {
	return &TicketCommentRepository{DB: db}
}

func (r *TicketCommentRepository) GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketComment, error) {
	var comments []*domain.TicketComment
	err := r.DB.Where("ticket_id = ?", ticketID).Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *TicketCommentRepository) Create(comment *domain.TicketComment) error {
	return r.DB.Create(comment).Error
}
