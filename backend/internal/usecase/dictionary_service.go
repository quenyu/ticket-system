package usecase

import (
	"ticket-system/backend/internal/model"
)

type DepartmentRepo interface {
	List() ([]interface {
		DepartmentID() int16
		DepartmentName() string
	}, error)
}

type TicketStatusRepo interface {
	List() ([]interface {
		ID() int16
		Code() string
		Label() string
	}, error)
}

type TicketPriorityRepo interface {
	List() ([]interface {
		ID() int16
		Code() string
		Label() string
		Level() int16
	}, error)
}

type DepartmentService struct {
	Repo DepartmentRepo
}

func NewDepartmentService(repo DepartmentRepo) *DepartmentService {
	return &DepartmentService{Repo: repo}
}

func (s *DepartmentService) List() ([]model.DepartmentDTO, error) {
	items, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	dtos := make([]model.DepartmentDTO, 0, len(items))
	for _, d := range items {
		dtos = append(dtos, model.DepartmentDTO{ID: d.DepartmentID(), Name: d.DepartmentName()})
	}
	return dtos, nil
}

type TicketStatusService struct {
	Repo TicketStatusRepo
}

func NewTicketStatusService(repo TicketStatusRepo) *TicketStatusService {
	return &TicketStatusService{Repo: repo}
}

func (s *TicketStatusService) List() ([]model.TicketStatusDTO, error) {
	items, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	dtos := make([]model.TicketStatusDTO, 0, len(items))
	for _, d := range items {
		dtos = append(dtos, model.TicketStatusDTO{ID: d.ID(), Code: d.Code(), Label: d.Label()})
	}
	return dtos, nil
}

type TicketPriorityService struct {
	Repo TicketPriorityRepo
}

func NewTicketPriorityService(repo TicketPriorityRepo) *TicketPriorityService {
	return &TicketPriorityService{Repo: repo}
}

func (s *TicketPriorityService) List() ([]model.TicketPriorityDTO, error) {
	items, err := s.Repo.List()
	if err != nil {
		return nil, err
	}
	dtos := make([]model.TicketPriorityDTO, 0, len(items))
	for _, d := range items {
		dtos = append(dtos, model.TicketPriorityDTO{ID: d.ID(), Code: d.Code(), Label: d.Label(), Level: d.Level()})
	}
	return dtos, nil
}
