-- Formdrain initial schema (STRICT tables).
-- Apply out-of-band before starting the app, e.g.:
--   sqlite3 /path/to/file.db < migrations/001_init.sql
--   DBURL=sqlite:/path/to/file.db
-- The HTTP process never runs DDL.
--
-- Note: SQLite STRICT only allows INTEGER/REAL/TEXT/BLOB/ANY. JSON is not a
-- STRICT datatype name; store JSON as TEXT and enforce with json_valid().

CREATE TABLE sites (
  id INTEGER PRIMARY KEY,
  hostname TEXT NOT NULL
) STRICT;

CREATE TABLE forms (
  id INTEGER PRIMARY KEY,
  site_id INTEGER NOT NULL,
  captcha_type TEXT NOT NULL,
  redirect_to TEXT,
  field_schema TEXT NOT NULL CHECK (json_valid(field_schema)),
  schema_version INTEGER NOT NULL,
  notifiers TEXT NOT NULL CHECK (json_valid(notifiers))
) STRICT;

CREATE TABLE form_responses (
  id TEXT PRIMARY KEY,
  form_id INTEGER NOT NULL,
  created_at TEXT NOT NULL,
  schema_version INTEGER NOT NULL,
  client_ip TEXT,
  payload TEXT NOT NULL CHECK (json_valid(payload))
) STRICT;

CREATE INDEX form_responses_form_id_created_at
  ON form_responses (form_id, created_at DESC);
