package repository

import (
	"ticket-system/backend/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("username = ? AND deleted_at IS NULL", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) List() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.DB.Where("deleted_at IS NULL").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
