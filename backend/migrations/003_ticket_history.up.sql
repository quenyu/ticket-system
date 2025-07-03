CREATE TABLE ticket_history (
    history_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL,
    ticket_created_at TIMESTAMPTZ NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    changed_by UUID NOT NULL REFERENCES users(user_id),
    field_name TEXT NOT NULL,
    old_value TEXT,
    new_value TEXT,
    change_data JSONB,
    FOREIGN KEY (ticket_id, ticket_created_at) REFERENCES tickets(ticket_id, created_at)
);