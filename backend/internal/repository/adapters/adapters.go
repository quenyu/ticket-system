package adapters

import (
	"ticket-system/backend/internal/domain"
	"ticket-system/backend/internal/repository"
)

type TicketStatusAdapter struct {
	Repo *repository.TicketStatusRepository
}
type ticketStatusWrap struct{ s *domain.TicketStatus }

func (w *ticketStatusWrap) ID() int16     { return w.s.ID }
func (w *ticketStatusWrap) Code() string  { return w.s.Code }
func (w *ticketStatusWrap) Label() string { return w.s.Label }
func (a *TicketStatusAdapter) List() ([]interface {
	ID() int16
	Code() string
	Label() string
}, error) {
	items, err := a.Repo.List()
	if err != nil {
		return nil, err
	}
	res := make([]interface {
		ID() int16
		Code() string
		Label() string
	}, len(items))
	for i, s := range items {
		res[i] = &ticketStatusWrap{s}
	}
	return res, nil
}

type TicketPriorityAdapter struct {
	Repo *repository.TicketPriorityRepository
}
type ticketPriorityWrap struct{ p *domain.TicketPriority }

func (w *ticketPriorityWrap) ID() int16     { return w.p.ID }
func (w *ticketPriorityWrap) Code() string  { return w.p.Code }
func (w *ticketPriorityWrap) Label() string { return w.p.Label }
func (w *ticketPriorityWrap) Level() int16  { return w.p.Level }
func (a *TicketPriorityAdapter) List() ([]interface {
	ID() int16
	Code() string
	Label() string
	Level() int16
}, error) {
	items, err := a.Repo.List()
	if err != nil {
		return nil, err
	}
	res := make([]interface {
		ID() int16
		Code() string
		Label() string
		Level() int16
	}, len(items))
	for i, p := range items {
		res[i] = &ticketPriorityWrap{p}
	}
	return res, nil
}

type DepartmentAdapter struct {
	Repo *repository.DepartmentRepository
}
type departmentWrap struct{ d *domain.Department }

func (w *departmentWrap) DepartmentID() int16    { return w.d.ID }
func (w *departmentWrap) DepartmentName() string { return w.d.Name }
func (a *DepartmentAdapter) List() ([]interface {
	DepartmentID() int16
	DepartmentName() string
}, error) {
	items, err := a.Repo.List()
	if err != nil {
		return nil, err
	}
	res := make([]interface {
		DepartmentID() int16
		DepartmentName() string
	}, len(items))
	for i, d := range items {
		dDept, ok := d.(*domain.Department)
		if !ok {
			panic("DepartmentAdapter: expected *domain.Department")
		}
		res[i] = &departmentWrap{dDept}
	}
	return res, nil
}

func NewTicketStatusAdapter(repo *repository.TicketStatusRepository) *TicketStatusAdapter {
	return &TicketStatusAdapter{Repo: repo}
}

func NewTicketPriorityAdapter(repo *repository.TicketPriorityRepository) *TicketPriorityAdapter {
	return &TicketPriorityAdapter{Repo: repo}
}

func NewDepartmentAdapter(repo *repository.DepartmentRepository) *DepartmentAdapter {
	return &DepartmentAdapter{Repo: repo}
}
