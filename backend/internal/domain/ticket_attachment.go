package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TicketAttachment struct {
	ID              uuid.UUID       `gorm:"column:attachment_id;primaryKey" json:"attachment_id"`
	TicketID        uuid.UUID       `gorm:"column:ticket_id" json:"ticket_id"`
	TicketCreatedAt time.Time       `gorm:"column:ticket_created_at" json:"ticket_created_at"`
	Filename        string          `gorm:"column:filename" json:"filename"`
	Metadata        json.RawMessage `gorm:"column:metadata" json:"metadata,omitempty"`
	FileData        []byte          `gorm:"column:file_data" json:"file_data,omitempty"`
	FilePath        *string         `gorm:"column:file_path" json:"file_path,omitempty"`
	UploadedBy      uuid.UUID       `gorm:"column:uploaded_by" json:"uploaded_by"`
	UploadedAt      time.Time       `gorm:"column:uploaded_at" json:"uploaded_at"`
}
