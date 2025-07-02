CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE DOMAIN email_address AS CITEXT
  CHECK (VALUE ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

CREATE TABLE roles (
    role_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE permissions (
    permission_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(role_id),
    permission_id UUID NOT NULL REFERENCES permissions(permission_id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE departments (
    department_id SMALLSERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE ticket_statuses (
    status_id SMALLSERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    label TEXT NOT NULL
);

CREATE TABLE ticket_priorities (
    priority_id SMALLSERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    label TEXT NOT NULL,
    level SMALLINT NOT NULL CHECK (level >= 1 AND level <= 5)
);

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
    ticket_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT,
    status_id SMALLINT NOT NULL REFERENCES ticket_statuses(status_id),
    priority_id SMALLINT NOT NULL REFERENCES ticket_priorities(priority_id),
    creator_id UUID NOT NULL REFERENCES users(user_id),
    assignee_id UUID REFERENCES users(user_id),
    department_id SMALLINT NOT NULL REFERENCES departments(department_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
) PARTITION BY RANGE (created_at);

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

CREATE TABLE ticket_comments (
    comment_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticket_id UUID NOT NULL REFERENCES tickets(ticket_id),
    author_id UUID NOT NULL REFERENCES users(user_id),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

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

CREATE INDEX idx_tickets_status_id ON tickets(status_id);
CREATE INDEX idx_tickets_priority_id ON tickets(priority_id);
CREATE INDEX idx_tickets_assignee_id ON tickets(assignee_id);
CREATE INDEX idx_tickets_department_id ON tickets(department_id);
CREATE INDEX idx_tickets_created_at ON tickets(created_at);
CREATE INDEX idx_tickets_updated_at ON tickets(updated_at);
CREATE INDEX idx_tickets_my_open ON tickets(assignee_id, status_id) WHERE deleted_at IS NULL;

ALTER TABLE tickets ADD COLUMN search_vector tsvector;
CREATE INDEX idx_tickets_search_vector ON tickets USING GIN (search_vector);
CREATE FUNCTION tickets_search_vector_trigger() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := to_tsvector('russian', coalesce(NEW.title,'') || ' ' || coalesce(NEW.description,''));
  RETURN NEW;
END
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_tickets_search_vector
  BEFORE INSERT OR UPDATE ON tickets
  FOR EACH ROW EXECUTE FUNCTION tickets_search_vector_trigger();

CREATE OR REPLACE FUNCTION fn_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_update_timestamp
  BEFORE UPDATE ON tickets
  FOR EACH ROW EXECUTE FUNCTION fn_update_timestamp();

CREATE TABLE ddl_audit (
    event_id SERIAL PRIMARY KEY,
    event_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    username TEXT,
    object_type TEXT,
    object_name TEXT,
    command_tag TEXT,
    ddl_command TEXT
);
CREATE OR REPLACE FUNCTION fn_ddl_audit() RETURNS event_trigger AS $$
BEGIN
  INSERT INTO ddl_audit(username, object_type, object_name, command_tag, ddl_command)
  SELECT
    session_user,
    objtype,
    objname,
    tag,
    command
  FROM pg_event_trigger_ddl_commands();
END;
$$ LANGUAGE plpgsql;
CREATE EVENT TRIGGER trg_ddl_audit
  ON ddl_command_end
  EXECUTE FUNCTION fn_ddl_audit();

CREATE MATERIALIZED VIEW mv_departments AS SELECT * FROM departments;
CREATE MATERIALIZED VIEW mv_ticket_statuses AS SELECT * FROM ticket_statuses;
CREATE MATERIALIZED VIEW mv_ticket_priorities AS SELECT * FROM ticket_priorities; 