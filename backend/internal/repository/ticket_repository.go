package repository

import (
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/usecase"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

func (r *TicketRepository) Create(ticket *domain.Ticket) error {
	if ticket.ID == uuid.Nil {
		ticket.ID = uuid.New()
	}
	now := time.Now().UTC()
	if ticket.CreatedAt.IsZero() {
		ticket.CreatedAt = now
	}
	ticket.UpdatedAt = now
	return r.DB.Create(ticket).Error
}

func (r *TicketRepository) GetByID(id uuid.UUID) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.DB.Where("ticket_id = ? AND deleted_at IS NULL", id).First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) Update(ticket *domain.Ticket) error {
	ticket.UpdatedAt = time.Now().UTC()
	return r.DB.Model(&domain.Ticket{}).
		Where("ticket_id = ? AND deleted_at IS NULL", ticket.ID).
		Updates(ticket).Error
}

func (r *TicketRepository) Delete(id uuid.UUID) error {
	return r.DB.Model(&domain.Ticket{}).
		Where("ticket_id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now().UTC()).Error
}

func (r *TicketRepository) Search(filter usecase.TicketFilter) ([]*domain.Ticket, error) {
	var tickets []*domain.Ticket
	db := r.DB.Model(&domain.Ticket{}).Where("deleted_at IS NULL")
	if filter.StatusID != nil {
		db = db.Where("status_id = ?", *filter.StatusID)
	}
	if filter.AssigneeID != nil {
		db = db.Where("assignee_id = ?", *filter.AssigneeID)
	}
	if filter.DepartmentID != nil {
		db = db.Where("department_id = ?", *filter.DepartmentID)
	}
	if filter.Q != "" {
		db = db.Where("search_vector @@ plainto_tsquery('russian', ?)", filter.Q)
		db = db.Order("ts_rank(search_vector, plainto_tsquery('russian', '" + filter.Q + "')) DESC")
	}
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}
	if err := db.Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}
