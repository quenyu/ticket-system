package domain

type TicketStatus struct {
	ID    int16  `gorm:"column:status_id;primaryKey"`
	Code  string `gorm:"column:code"`
	Label string `gorm:"column:label"`
}
