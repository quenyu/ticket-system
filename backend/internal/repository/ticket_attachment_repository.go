package repository

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketAttachmentRepository struct {
	DB *gorm.DB
}

func NewTicketAttachmentRepository(db *gorm.DB) *TicketAttachmentRepository {
	return &TicketAttachmentRepository{DB: db}
}

func (r *TicketAttachmentRepository) GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketAttachment, error) {
	var atts []*domain.TicketAttachment
	err := r.DB.Where("ticket_id = ?", ticketID).Order("uploaded_at ASC").Find(&atts).Error
	if err != nil {
		return nil, err
	}
	return atts, nil
}

func (r *TicketAttachmentRepository) GetByID(attID uuid.UUID) (*domain.TicketAttachment, error) {
	var att domain.TicketAttachment
	err := r.DB.Where("attachment_id = ?", attID).First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *TicketAttachmentRepository) Create(att *domain.TicketAttachment) error {
	return r.DB.Create(att).Error
}

func (r *TicketAttachmentRepository) Delete(attID uuid.UUID) error {
	return r.DB.Delete(&domain.TicketAttachment{}, "attachment_id = ?", attID).Error
}
