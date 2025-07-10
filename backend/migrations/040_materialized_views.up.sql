CREATE MATERIALIZED VIEW mv_departments AS SELECT * FROM departments;
CREATE MATERIALIZED VIEW mv_ticket_statuses AS SELECT * FROM ticket_statuses;
CREATE MATERIALIZED VIEW mv_ticket_priorities AS SELECT * FROM ticket_priorities; 

INSERT INTO roles (role_id, name) VALUES
  ('00000000-0000-0000-0000-000000000002', 'Admin'),
  ('00000000-0000-0000-0000-000000000001', 'User');

INSERT INTO ticket_statuses (code, label) VALUES
  ('open', 'Open'),
  ('closed', 'Closed');

INSERT INTO ticket_priorities (code, label, level) VALUES
  ('low', 'Low', 1),
  ('medium', 'Medium', 3),
  ('high', 'High', 5);

INSERT INTO departments (name) VALUES
  ('IT'),
  ('HR');