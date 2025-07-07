CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username CITEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email email_address NOT NULL,
    role_id UUID NOT NULL REFERENCES roles(role_id),
    department_id SMALLINT NOT NULL REFERENCES departments(department_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE tickets (
    ticket_id UUID NOT NULL DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    status_id SMALLINT NOT NULL REFERENCES ticket_statuses(status_id),
    priority_id SMALLINT NOT NULL REFERENCES ticket_priorities(priority_id),
    creator_id UUID NOT NULL REFERENCES users(user_id),
    assignee_id UUID REFERENCES users(user_id),
    department_id SMALLINT NOT NULL REFERENCES departments(department_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (ticket_id, created_at)
) PARTITION BY RANGE (created_at);

CREATE TABLE tickets_2024 PARTITION OF tickets
    FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
CREATE TABLE tickets_2025 PARTITION OF tickets
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
CREATE TABLE tickets_2026 PARTITION OF tickets
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');