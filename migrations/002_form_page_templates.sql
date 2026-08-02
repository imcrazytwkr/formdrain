-- Optional owner Mustache pages for form submit HTML responses.
-- NULL means use the embedded system template for that page.
-- Apply out-of-band after 001_init.sql, e.g.:
--   sqlite3 /path/to/file.db < migrations/002_form_page_templates.sql

ALTER TABLE forms ADD COLUMN success_template TEXT;
ALTER TABLE forms ADD COLUMN error_template TEXT;
ALTER TABLE forms ADD COLUMN redirect_template TEXT;
