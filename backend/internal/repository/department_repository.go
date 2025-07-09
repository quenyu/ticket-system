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

func (r *DepartmentRepository) List() ([]interface {
	DepartmentID() int16
	DepartmentName() string
}, error) {
	var departments []*domain.Department
	if err := r.DB.Find(&departments).Error; err != nil {
		return nil, err
	}
	res := make([]interface {
		DepartmentID() int16
		DepartmentName() string
	}, len(departments))
	for i, d := range departments {
		res[i] = d
	}
	return res, nil
}
