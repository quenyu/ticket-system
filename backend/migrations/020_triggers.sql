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