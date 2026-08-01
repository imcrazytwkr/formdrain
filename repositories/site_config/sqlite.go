package site_config

import (
	"context"
	"database/sql"
	"errors"

	"github.com/imcrazytwkr/formdrain/models/site_config"
	"github.com/imcrazytwkr/formdrain/repositories"
)

type sqliteSiteConfigRepository struct {
	db *sql.DB
}

func NewSqliteSiteConfigRepository(db *sql.DB) repositories.SiteConfigRepository {
	return &sqliteSiteConfigRepository{db: db}
}

const selectSiteConfigById = `SELECT hostname FROM sites WHERE id = ?`

func (r *sqliteSiteConfigRepository) GetSiteConfigById(ctx context.Context, id int64) (*site_config.SiteConfig, error) {
	if id < 1 {
		return nil, nil
	}

	var config site_config.SiteConfig
	err := r.db.QueryRowContext(ctx, selectSiteConfigById, id).Scan(&config.Hostname)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	config.SiteId = id
	return &config, nil
}
