package usecase

import (
	"ticket-system/backend/internal/model"
)

type DepartmentService struct {
	Repo DepartmentRepository
}

func NewDepartmentService(repo DepartmentRepository) *DepartmentService {
	return &DepartmentService{Repo: repo}
}

func (s *DepartmentService) List() ([]model.DepartmentDTO, error) {
	items, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	dtos := make([]model.DepartmentDTO, 0, len(items))
	for _, d := range items {
		dtos = append(dtos, model.DepartmentDTO{ID: d.ID, Name: d.Name})
	}
	return dtos, nil
}

type TicketStatusService struct {
	Repo TicketStatusRepository
}

func NewTicketStatusService(repo TicketStatusRepository) *TicketStatusService {
	return &TicketStatusService{Repo: repo}
}

func (s *TicketStatusService) List() ([]model.TicketStatusDTO, error) {
	items, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	dtos := make([]model.TicketStatusDTO, 0, len(items))
	for _, d := range items {
		dtos = append(dtos, model.TicketStatusDTO{ID: d.ID, Code: d.Code, Label: d.Label})
	}
	return dtos, nil
}

type TicketPriorityService struct {
	Repo TicketPriorityRepository
}

func NewTicketPriorityService(repo TicketPriorityRepository) *TicketPriorityService {
	return &TicketPriorityService{Repo: repo}
}

func (s *TicketPriorityService) List() ([]model.TicketPriorityDTO, error) {
	items, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	dtos := make([]model.TicketPriorityDTO, 0, len(items))
	for _, d := range items {
		dtos = append(dtos, model.TicketPriorityDTO{ID: d.ID, Code: d.Code, Label: d.Label, Level: d.Level})
	}
	return dtos, nil
}
