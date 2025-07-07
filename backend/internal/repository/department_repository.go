package repository

import (
	"ticket-system/backend/internal/domain"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	DB *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{DB: db}
}

func (r *DepartmentRepository) List() ([]*domain.Department, error) {
	var departments []*domain.Department
	if err := r.DB.Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}
