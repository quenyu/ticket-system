package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TicketHistory struct {
	ID         uuid.UUID
	TicketID   uuid.UUID
	ChangedAt  time.Time
	ChangedBy  uuid.UUID
	FieldName  string
	OldValue   string
	NewValue   string
	ChangeData json.RawMessage
}
