package repositories

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/site_config"
)

type SiteConfigRepository interface {
	GetSiteConfigById(ctx context.Context, id int64) (*site_config.SiteConfig, error)
	ListByOwnerIDAfterID(ctx context.Context, ownerID, afterID int64, limit int) ([]*site_config.SiteConfig, error)
	ListByOwnerIDAfterHostname(ctx context.Context, ownerID int64, afterHostname string, afterID int64, limit int) ([]*site_config.SiteConfig, error)
}
