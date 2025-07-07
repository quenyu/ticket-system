package usecase

import "ticket-system/backend/internal/domain"

type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	Create(user *domain.User) error
}
