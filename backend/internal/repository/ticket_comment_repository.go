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

func (r *TicketCommentRepository) Delete(commentID uuid.UUID) error {
	return r.DB.Delete(&domain.TicketComment{}, "comment_id = ?", commentID).Error
}

func (r *TicketCommentRepository) Update(comment *domain.TicketComment) error {
	return r.DB.Model(&domain.TicketComment{}).
		Where("comment_id = ?", comment.ID).
		Updates(map[string]interface{}{
			"content":    comment.Content,
			"updated_at": comment.CreatedAt,
		}).Error
}

func (r *TicketCommentRepository) GetByID(commentID uuid.UUID) (*domain.TicketComment, error) {
	var comment domain.TicketComment
	if err := r.DB.Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}
