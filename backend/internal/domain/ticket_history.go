package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TicketHistory struct {
	ID              uuid.UUID       `gorm:"column:history_id;primaryKey" json:"history_id"`
	TicketID        uuid.UUID       `gorm:"column:ticket_id" json:"ticket_id"`
	TicketCreatedAt time.Time       `gorm:"column:ticket_created_at" json:"ticket_created_at"`
	ChangedAt       time.Time       `gorm:"column:changed_at" json:"changed_at"`
	ChangedBy       uuid.UUID       `gorm:"column:changed_by" json:"changed_by"`
	FieldName       string          `gorm:"column:field_name" json:"field_name"`
	OldValue        string          `gorm:"column:old_value" json:"old_value"`
	NewValue        string          `gorm:"column:new_value" json:"new_value"`
	ChangeData      json.RawMessage `gorm:"column:change_data" json:"change_data,omitempty"`
}
