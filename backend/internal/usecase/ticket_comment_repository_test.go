package usecase

import (
	"errors"
	"testing"

	"ticket-system/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockTicketCommentRepo struct {
	GetByTicketIDFunc func(ticketID uuid.UUID) ([]*domain.TicketComment, error)
	GetByIDFunc       func(commentID uuid.UUID) (*domain.TicketComment, error)
	CreateFunc        func(comment *domain.TicketComment) error
	DeleteFunc        func(commentID uuid.UUID) error
}

func (m *mockTicketCommentRepo) GetByTicketID(ticketID uuid.UUID) ([]*domain.TicketComment, error) {
	return m.GetByTicketIDFunc(ticketID)
}
func (m *mockTicketCommentRepo) GetByID(commentID uuid.UUID) (*domain.TicketComment, error) {
	return m.GetByIDFunc(commentID)
}
func (m *mockTicketCommentRepo) Create(comment *domain.TicketComment) error {
	return m.CreateFunc(comment)
}
func (m *mockTicketCommentRepo) Delete(commentID uuid.UUID) error {
	return m.DeleteFunc(commentID)
}

func TestTicketCommentRepository_GetByTicketID_OK(t *testing.T) {
	repo := &mockTicketCommentRepo{
		GetByTicketIDFunc: func(ticketID uuid.UUID) ([]*domain.TicketComment, error) {
			return []*domain.TicketComment{{TicketID: ticketID}}, nil
		},
	}
	ticketID := uuid.New()
	comments, err := repo.GetByTicketID(ticketID)
	assert.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, ticketID, comments[0].TicketID)
}

func TestTicketCommentRepository_GetByTicketID_Error(t *testing.T) {
	repo := &mockTicketCommentRepo{
		GetByTicketIDFunc: func(ticketID uuid.UUID) ([]*domain.TicketComment, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := repo.GetByTicketID(uuid.New())
	assert.Error(t, err)
}

func TestTicketCommentRepository_GetByID_OK(t *testing.T) {
	repo := &mockTicketCommentRepo{
		GetByIDFunc: func(commentID uuid.UUID) (*domain.TicketComment, error) {
			return &domain.TicketComment{ID: commentID}, nil
		},
	}
	commentID := uuid.New()
	comment, err := repo.GetByID(commentID)
	assert.NoError(t, err)
	assert.Equal(t, commentID, comment.ID)
}

func TestTicketCommentRepository_GetByID_Error(t *testing.T) {
	repo := &mockTicketCommentRepo{
		GetByIDFunc: func(commentID uuid.UUID) (*domain.TicketComment, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := repo.GetByID(uuid.New())
	assert.Error(t, err)
}

func TestTicketCommentRepository_Create_OK(t *testing.T) {
	repo := &mockTicketCommentRepo{
		CreateFunc: func(comment *domain.TicketComment) error { return nil },
	}
	assert.NoError(t, repo.Create(&domain.TicketComment{}))
}

func TestTicketCommentRepository_Create_Error(t *testing.T) {
	repo := &mockTicketCommentRepo{
		CreateFunc: func(comment *domain.TicketComment) error { return errors.New("fail") },
	}
	assert.Error(t, repo.Create(&domain.TicketComment{}))
}

func TestTicketCommentRepository_Delete_OK(t *testing.T) {
	repo := &mockTicketCommentRepo{
		DeleteFunc: func(commentID uuid.UUID) error { return nil },
	}
	assert.NoError(t, repo.Delete(uuid.New()))
}

func TestTicketCommentRepository_Delete_Error(t *testing.T) {
	repo := &mockTicketCommentRepo{
		DeleteFunc: func(commentID uuid.UUID) error { return errors.New("fail") },
	}
	assert.Error(t, repo.Delete(uuid.New()))
}
