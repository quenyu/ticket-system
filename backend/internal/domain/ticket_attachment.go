package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TicketAttachment struct {
	ID         uuid.UUID
	TicketID   uuid.UUID
	Filename   string
	Metadata   json.RawMessage
	FileData   []byte
	FilePath   string
	UploadedBy uuid.UUID
	UploadedAt time.Time
}
