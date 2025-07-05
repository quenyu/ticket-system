package usecase

import (
	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
)

type TicketAttachmentRepository interface {
	GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketAttachment, error)
	GetByID(attID uuid.UUID) (*domain.TicketAttachment, error)
	Create(att *domain.TicketAttachment) error
	Delete(attID uuid.UUID) error
}
