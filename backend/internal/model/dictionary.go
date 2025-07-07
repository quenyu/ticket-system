package model

type DepartmentDTO struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type TicketStatusDTO struct {
	ID    int16  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
}

type TicketPriorityDTO struct {
	ID    int16  `json:"id"`
	Code  string `json:"code"`
	Label string `json:"label"`
	Level int16  `json:"level"`
}
