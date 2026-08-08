-- Formdrain initial schema (STRICT tables).
-- Apply out-of-band before starting the app, e.g.:
--   sqlite3 /path/to/file.db < migrations/001_init.sql
--   DBURL=sqlite:/path/to/file.db
-- The HTTP process never runs DDL.
--
-- Note: SQLite STRICT only allows INTEGER/REAL/TEXT/BLOB/ANY. JSON is not a
-- STRICT datatype name; store JSON as TEXT and enforce with json_valid().
-- Runtime connections must enable foreign keys: PRAGMA foreign_keys = ON;

CREATE TABLE accounts (
  id INTEGER PRIMARY KEY,
  email TEXT NOT NULL COLLATE NOCASE,
  password_hash TEXT NOT NULL,
  created_at TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX unique_account_email ON accounts (email);

CREATE TABLE sites (
  id INTEGER PRIMARY KEY,
  hostname TEXT NOT NULL,
  owner_id INTEGER NOT NULL REFERENCES accounts (id)
) STRICT;

CREATE TABLE forms (
  id INTEGER PRIMARY KEY,
  site_id INTEGER NOT NULL REFERENCES sites (id),
  captcha_type TEXT NOT NULL,
  redirect_to TEXT,
  field_schema TEXT NOT NULL CHECK (json_valid(field_schema)),
  schema_version INTEGER NOT NULL,
  notifiers TEXT NOT NULL CHECK (json_valid(notifiers)),
  captcha_field TEXT
) STRICT;

CREATE TABLE form_responses (
  id TEXT PRIMARY KEY,
  client_ip TEXT,
  form_id INTEGER NOT NULL REFERENCES forms (id),
  schema_version INTEGER NOT NULL,
  payload TEXT NOT NULL CHECK (json_valid(payload)),
  created_at TEXT NOT NULL
) STRICT;

CREATE INDEX form_responses_form_id_created_at
  ON form_responses (form_id, created_at DESC);
