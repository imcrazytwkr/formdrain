-- Optional form field name for the captcha response token.
-- NULL/empty means use the provider default (h-captcha / g-recaptcha).
-- Apply out-of-band after 001_init.sql, e.g.:
--   sqlite3 /path/to/file.db < migrations/002_form_captcha_field.sql

ALTER TABLE forms ADD COLUMN captcha_field TEXT;
