package usecase

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/model"
)

type mockUserRepo struct {
	FindByUsernameFunc func(username string) (*domain.User, error)
	CreateFunc         func(user *domain.User) error
}

func (m *mockUserRepo) FindByUsername(username string) (*domain.User, error) {
	return m.FindByUsernameFunc(username)
}
func (m *mockUserRepo) Create(user *domain.User) error {
	return m.CreateFunc(user)
}

func TestAuthService_Login_Success(t *testing.T) {
	password := "secret"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: string(hash),
		RoleID:       uuid.New(),
		DepartmentID: 1,
	}
	repo := &mockUserRepo{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return user, nil
		},
	}
	svc := NewAuthService(repo, "testsecret")
	resp, err := svc.Login("testuser", password)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, resp.User.Username)
	assert.NotEmpty(t, resp.Token)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("right"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: string(hash),
		RoleID:       uuid.New(),
		DepartmentID: 1,
	}
	repo := &mockUserRepo{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return user, nil
		},
	}
	svc := NewAuthService(repo, "testsecret")
	_, err := svc.Login("testuser", "wrong")
	assert.Error(t, err)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{
		FindByUsernameFunc: func(username string) (*domain.User, error) {
			return nil, errors.New("not found")
		},
	}
	svc := NewAuthService(repo, "testsecret")
	_, err := svc.Login("nouser", "pass")
	assert.Error(t, err)
}

func TestAuthService_Register_Success(t *testing.T) {
	var created *domain.User
	repo := &mockUserRepo{
		CreateFunc: func(user *domain.User) error {
			created = user
			return nil
		},
	}
	svc := NewAuthService(repo, "testsecret")
	input := model.UserRegisterRequest{
		Username:     "newuser",
		Password:     "pass123",
		Email:        "new@ex.com",
		RoleID:       uuid.New().String(),
		DepartmentID: 2,
	}
	resp, err := svc.Register(input)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.UserID)
	assert.Equal(t, input.Username, created.Username)
	assert.Equal(t, input.Email, created.Email)
	assert.Equal(t, input.DepartmentID, created.DepartmentID)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(created.PasswordHash), []byte(input.Password)))
}
