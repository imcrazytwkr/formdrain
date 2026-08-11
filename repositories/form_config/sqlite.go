package form_config

import (
	"database/sql"

	"github.com/imcrazytwkr/formdrain/repositories"
)

type sqliteFormConfigRepository struct {
	db *sql.DB
}

func NewSqliteFormConfigRepository(db *sql.DB) repositories.FormConfigRepository {
	return &sqliteFormConfigRepository{db: db}
}
