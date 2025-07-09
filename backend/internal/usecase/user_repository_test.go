package usecase

import (
	"errors"
	"testing"

	"ticket-system/backend/internal/domain"

	"github.com/stretchr/testify/assert"
)

type mockUserRepoUC struct {
	FindByUsernameFunc func(username string) (*domain.User, error)
	CreateFunc         func(user *domain.User) error
}

func (m *mockUserRepoUC) FindByUsername(username string) (*domain.User, error) {
	return m.FindByUsernameFunc(username)
}
func (m *mockUserRepoUC) Create(user *domain.User) error {
	return m.CreateFunc(user)
}

func TestUserRepository_FindByUsername_OK(t *testing.T) {
	repo := &mockUserRepoUC{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return &domain.User{Username: username}, nil
		},
	}
	user, err := repo.FindByUsername("alice")
	assert.NoError(t, err)
	assert.Equal(t, "alice", user.Username)
}

func TestUserRepository_FindByUsername_Error(t *testing.T) {
	repo := &mockUserRepoUC{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return nil, errors.New("fail")
		},
	}
	_, err := repo.FindByUsername("bob")
	assert.Error(t, err)
}

func TestUserRepository_Create_OK(t *testing.T) {
	repo := &mockUserRepoUC{
		CreateFunc: func(user *domain.User) error { return nil },
	}
	assert.NoError(t, repo.Create(&domain.User{}))
}

func TestUserRepository_Create_Error(t *testing.T) {
	repo := &mockUserRepoUC{
		CreateFunc: func(user *domain.User) error { return errors.New("fail") },
	}
	assert.Error(t, repo.Create(&domain.User{}))
}
