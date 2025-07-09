package usecase

import (
	"errors"
	"testing"

	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockTicketRepo struct {
	CreateFunc  func(ticket *domain.Ticket) error
	GetByIDFunc func(id uuid.UUID) (*domain.Ticket, error)
	UpdateFunc  func(ticket *domain.Ticket) error
	DeleteFunc  func(id uuid.UUID) error
	SearchFunc  func(filter model.TicketFilter) ([]*domain.Ticket, error)
}

func (m *mockTicketRepo) Create(ticket *domain.Ticket) error           { return m.CreateFunc(ticket) }
func (m *mockTicketRepo) GetByID(id uuid.UUID) (*domain.Ticket, error) { return m.GetByIDFunc(id) }
func (m *mockTicketRepo) Update(ticket *domain.Ticket) error           { return m.UpdateFunc(ticket) }
func (m *mockTicketRepo) Delete(id uuid.UUID) error                    { return m.DeleteFunc(id) }
func (m *mockTicketRepo) Search(filter model.TicketFilter) ([]*domain.Ticket, error) {
	return m.SearchFunc(filter)
}

func TestTicketRepository_Create_OK(t *testing.T) {
	repo := &mockTicketRepo{
		CreateFunc: func(ticket *domain.Ticket) error { return nil },
	}
	err := repo.Create(&domain.Ticket{})
	assert.NoError(t, err)
}

func TestTicketRepository_Create_Error(t *testing.T) {
	repo := &mockTicketRepo{
		CreateFunc: func(ticket *domain.Ticket) error { return errors.New("fail") },
	}
	err := repo.Create(&domain.Ticket{})
	assert.Error(t, err)
}

func TestTicketRepository_GetByID_OK(t *testing.T) {
	id := uuid.New()
	repo := &mockTicketRepo{
		GetByIDFunc: func(got uuid.UUID) (*domain.Ticket, error) {
			assert.Equal(t, id, got)
			return &domain.Ticket{ID: id}, nil
		},
	}
	ticket, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, id, ticket.ID)
}

func TestTicketRepository_GetByID_Error(t *testing.T) {
	repo := &mockTicketRepo{
		GetByIDFunc: func(id uuid.UUID) (*domain.Ticket, error) { return nil, errors.New("fail") },
	}
	_, err := repo.GetByID(uuid.New())
	assert.Error(t, err)
}

func TestTicketRepository_Update_OK(t *testing.T) {
	repo := &mockTicketRepo{
		UpdateFunc: func(ticket *domain.Ticket) error { return nil },
	}
	assert.NoError(t, repo.Update(&domain.Ticket{}))
}

func TestTicketRepository_Update_Error(t *testing.T) {
	repo := &mockTicketRepo{
		UpdateFunc: func(ticket *domain.Ticket) error { return errors.New("fail") },
	}
	assert.Error(t, repo.Update(&domain.Ticket{}))
}

func TestTicketRepository_Delete_OK(t *testing.T) {
	repo := &mockTicketRepo{
		DeleteFunc: func(id uuid.UUID) error { return nil },
	}
	assert.NoError(t, repo.Delete(uuid.New()))
}

func TestTicketRepository_Delete_Error(t *testing.T) {
	repo := &mockTicketRepo{
		DeleteFunc: func(id uuid.UUID) error { return errors.New("fail") },
	}
	assert.Error(t, repo.Delete(uuid.New()))
}

func TestTicketRepository_Search_OK(t *testing.T) {
	repo := &mockTicketRepo{
		SearchFunc: func(filter model.TicketFilter) ([]*domain.Ticket, error) {
			return []*domain.Ticket{{ID: uuid.New()}}, nil
		},
	}
	ts, err := repo.Search(model.TicketFilter{})
	assert.NoError(t, err)
	assert.Len(t, ts, 1)
}

func TestTicketRepository_Search_Error(t *testing.T) {
	repo := &mockTicketRepo{
		SearchFunc: func(filter model.TicketFilter) ([]*domain.Ticket, error) { return nil, errors.New("fail") },
	}
	_, err := repo.Search(model.TicketFilter{})
	assert.Error(t, err)
}
