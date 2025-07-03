package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TicketAttachment struct {
	ID              uuid.UUID       `gorm:"column:attachment_id;primaryKey"`
	TicketID        uuid.UUID       `gorm:"column:ticket_id"`
	TicketCreatedAt time.Time       `gorm:"column:ticket_created_at"`
	Filename        string          `gorm:"column:filename"`
	Metadata        json.RawMessage `gorm:"column:metadata"`
	FileData        []byte          `gorm:"column:file_data"`
	FilePath        string          `gorm:"column:file_path"`
	UploadedBy      uuid.UUID       `gorm:"column:uploaded_by"`
	UploadedAt      time.Time       `gorm:"column:uploaded_at"`
}
