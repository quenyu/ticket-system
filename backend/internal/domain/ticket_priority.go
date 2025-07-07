package domain

type TicketPriority struct {
	ID    int16  `gorm:"column:priority_id;primaryKey"`
	Code  string `gorm:"column:code"`
	Label string `gorm:"column:label"`
	Level int16  `gorm:"column:level"`
}
