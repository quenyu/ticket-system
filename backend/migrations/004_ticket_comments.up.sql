CREATE TABLE ticket_comments (
    comment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL,
    ticket_created_at TIMESTAMPTZ NOT NULL,
    author_id UUID NOT NULL REFERENCES users(user_id),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (ticket_id, ticket_created_at) REFERENCES tickets(ticket_id, created_at)
);