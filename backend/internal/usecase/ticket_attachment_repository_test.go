package usecase

import (
	"errors"
	"testing"

	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockTicketAttachmentRepo struct {
	GetByTicketIDFunc func(ticketID uuid.UUID) ([]*domain.TicketAttachment, error)
	GetByIDFunc       func(attID uuid.UUID) (*domain.TicketAttachment, error)
	CreateFunc        func(att *domain.TicketAttachment) error
	DeleteFunc        func(attID uuid.UUID) error
}

func (m *mockTicketAttachmentRepo) GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketAttachment, error) {
	return m.GetByTicketIDFunc(ticketID)
}
func (m *mockTicketAttachmentRepo) GetByID(attID uuid.UUID) (*domain.TicketAttachment, error) {
	return m.GetByIDFunc(attID)
}
func (m *mockTicketAttachmentRepo) Create(att *domain.TicketAttachment) error {
	return m.CreateFunc(att)
}
func (m *mockTicketAttachmentRepo) Delete(attID uuid.UUID) error {
	return m.DeleteFunc(attID)
}

func TestTicketAttachmentRepository_GetByTicketID_OK(t *testing.T) {
	repo := &mockTicketAttachmentRepo{
		GetByTicketIDFunc: func(ticketID uuid.UUID) ([]*domain.TicketAttachment, error) {
			return []*domain.TicketAttachment{{TicketID: ticketID}}, nil
		},
	}
	ticketID := uuid.New()
	atts, err := repo.GetByTicketID(ticketID)
	assert.NoError(t, err)
	assert.Len(t, atts, 1)
	assert.Equal(t, ticketID, atts[0].TicketID)
}

func TestTicketAttachmentRepository_GetByTicketID_Error(t *testing.T) {
	repo := &mockTicketAttachmentRepo{
		GetByTicketIDFunc: func(ticketID uuid.UUID) ([]*domain.TicketAttachment, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := repo.GetByTicketID(uuid.New())
	assert.Error(t, err)
}

func TestTicketAttachmentRepository_GetByID_OK(t *testing.T) {
	repo := &mockTicketAttachmentRepo{
		GetByIDFunc: func(attID uuid.UUID) (*domain.TicketAttachment, error) {
			return &domain.TicketAttachment{ID: attID}, nil
		},
	}
	attID := uuid.New()
	att, err := repo.GetByID(attID)
	assert.NoError(t, err)
	assert.Equal(t, attID, att.ID)
}

func TestTicketAttachmentRepository_GetByID_Error(t *testing.T) {
	repo := &mockTicketAttachmentRepo{
		GetByIDFunc: func(attID uuid.UUID) (*domain.TicketAttachment, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := repo.GetByID(uuid.New())
	assert.Error(t, err)
}

func TestTicketAttachmentRepository_Create_OK(t *testing.T) {
	repo := &mockTicketAttachmentRepo{
		CreateFunc: func(att *domain.TicketAttachment) error { return nil },
	}
	assert.NoError(t, repo.Create(&domain.TicketAttachment{}))
}

func TestTicketAttachmentRepository_Create_Error(t *testing.T) {
	repo := &mockTicketAttachmentRepo{
		CreateFunc: func(att *domain.TicketAttachment) error { return errors.New("fail") },
	}
	assert.Error(t, repo.Create(&domain.TicketAttachment{}))
}

func TestTicketAttachmentRepository_Delete_OK(t *testing.T) {
	repo := &mockTicketAttachmentRepo{
		DeleteFunc: func(attID uuid.UUID) error { return nil },
	}
	assert.NoError(t, repo.Delete(uuid.New()))
}

func TestTicketAttachmentRepository_Delete_Error(t *testing.T) {
	repo := &mockTicketAttachmentRepo{
		DeleteFunc: func(attID uuid.UUID) error { return errors.New("fail") },
	}
	assert.Error(t, repo.Delete(uuid.New()))
}
