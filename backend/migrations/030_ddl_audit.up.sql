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
DECLARE
  rec record;
BEGIN
  FOR rec IN SELECT * FROM pg_event_trigger_ddl_commands() LOOP
    INSERT INTO ddl_audit(username, object_type, object_name, command_tag, ddl_command)
    VALUES (
      session_user,
      rec.object_type,
      rec.object_identity,
      rec.command_tag,
      NULL
    );
  END LOOP;
END;
$$ LANGUAGE plpgsql;

CREATE EVENT TRIGGER trg_ddl_audit
  ON ddl_command_end
  EXECUTE FUNCTION fn_ddl_audit();