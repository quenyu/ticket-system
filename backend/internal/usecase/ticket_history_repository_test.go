package usecase

import (
	"errors"
	"testing"

	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockTicketHistoryRepo struct {
	GetByTicketIDFunc func(ticketID uuid.UUID) ([]*domain.TicketHistory, error)
}

func (m *mockTicketHistoryRepo) GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketHistory, error) {
	return m.GetByTicketIDFunc(ticketID)
}

func TestTicketHistoryRepository_GetByTicketID_OK(t *testing.T) {
	repo := &mockTicketHistoryRepo{
		GetByTicketIDFunc: func(ticketID uuid.UUID) ([]*domain.TicketHistory, error) {
			return []*domain.TicketHistory{{TicketID: ticketID}}, nil
		},
	}
	ticketID := uuid.New()
	hist, err := repo.GetByTicketID(ticketID)
	assert.NoError(t, err)
	assert.Len(t, hist, 1)
	assert.Equal(t, ticketID, hist[0].TicketID)
}

func TestTicketHistoryRepository_GetByTicketID_Error(t *testing.T) {
	repo := &mockTicketHistoryRepo{
		GetByTicketIDFunc: func(ticketID uuid.UUID) ([]*domain.TicketHistory, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := repo.GetByTicketID(uuid.New())
	assert.Error(t, err)
}
