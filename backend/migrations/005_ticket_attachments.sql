CREATE TABLE ticket_attachments (
    attachment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(ticket_id),
    filename TEXT NOT NULL,
    metadata JSONB,
    file_data BYTEA,
    file_path TEXT,
    uploaded_by UUID NOT NULL REFERENCES users(user_id),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (
        (file_data IS NOT NULL AND file_path IS NULL)
        OR (file_data IS NULL AND file_path IS NOT NULL)
    )
); 