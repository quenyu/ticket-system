package usecase

import "ticket-system/backend/internal/domain"

type DepartmentRepository interface {
	List() ([]*domain.Department, error)
}

type TicketStatusRepository interface {
	List() ([]*domain.TicketStatus, error)
}

type TicketPriorityRepository interface {
	List() ([]*domain.TicketPriority, error)
}
