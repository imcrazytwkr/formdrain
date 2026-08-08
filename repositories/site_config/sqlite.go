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

const selectSiteConfigById = `SELECT hostname, owner_id FROM sites WHERE id = ?`

func (r *sqliteSiteConfigRepository) GetSiteConfigById(ctx context.Context, id int64) (*site_config.SiteConfig, error) {
	if id < 1 {
		return nil, nil
	}

	var config site_config.SiteConfig
	err := r.db.QueryRowContext(ctx, selectSiteConfigById, id).Scan(&config.Hostname, &config.OwnerId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	config.SiteId = id
	return &config, nil
}

const listCommon = `
SELECT id, hostname, owner_id
FROM sites
`

const listByOwnerAfterID = listCommon + `
WHERE owner_id = ? AND id > ?
ORDER BY id ASC
LIMIT ?`

func (r *sqliteSiteConfigRepository) ListByOwnerIDAfterID(ctx context.Context, ownerID, afterID int64, limit int) ([]*site_config.SiteConfig, error) {
	if ownerID < 1 || limit < 1 {
		return []*site_config.SiteConfig{}, nil
	}

	rows, err := r.db.QueryContext(ctx, listByOwnerAfterID, ownerID, afterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSiteConfigs(rows)
}

const listByOwnerHostname = listCommon + `
WHERE owner_id = ?
`

const listByOwnerHostnameOrder = `
ORDER BY hostname ASC, id ASC
LIMIT ?
`

const listByOwnerHostnameFirstPage = listByOwnerHostname + listByOwnerHostnameOrder

const listByOwnerAfterHostname = listByOwnerHostname + `
  AND (hostname > ? OR (hostname = ? AND id > ?))
` + listByOwnerHostnameOrder

func (r *sqliteSiteConfigRepository) ListByOwnerIDAfterHostname(ctx context.Context, ownerID int64, afterHostname string, afterID int64, limit int) ([]*site_config.SiteConfig, error) {
	if ownerID < 1 || limit < 1 {
		return []*site_config.SiteConfig{}, nil
	}

	var rows *sql.Rows
	var err error
	if len(afterHostname) < 1 && afterID == 0 {
		rows, err = r.db.QueryContext(ctx, listByOwnerHostnameFirstPage, ownerID, limit)
	} else {
		rows, err = r.db.QueryContext(ctx, listByOwnerAfterHostname, ownerID, afterHostname, afterHostname, afterID, limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanSiteConfigs(rows)
}

func scanSiteConfigs(rows *sql.Rows) ([]*site_config.SiteConfig, error) {
	out := make([]*site_config.SiteConfig, 0)
	for rows.Next() {
		var cfg site_config.SiteConfig
		if err := rows.Scan(&cfg.SiteId, &cfg.Hostname, &cfg.OwnerId); err != nil {
			return nil, err
		}
		out = append(out, &cfg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
