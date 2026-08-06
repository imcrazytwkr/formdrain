-- Accounts for login (STRICT tables).
-- Apply out-of-band before starting the app, e.g.:
--   sqlite3 /path/to/file.db < migrations/003_accounts.sql
-- The HTTP process never runs DDL.

CREATE TABLE accounts (
  id INTEGER PRIMARY KEY,
  email TEXT NOT NULL COLLATE NOCASE,
  password_hash TEXT NOT NULL,
  created_at TEXT NOT NULL
) STRICT;

CREATE UNIQUE INDEX accounts_email ON accounts (email);
