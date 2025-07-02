CREATE TABLE ticket_history (
    history_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(ticket_id),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    changed_by UUID NOT NULL REFERENCES users(user_id),
    field_name TEXT NOT NULL,
    old_value TEXT,
    new_value TEXT,
    change_data JSONB
); 