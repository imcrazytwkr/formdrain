package repositories

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/site_config"
)

type SiteConfigRepository interface {
	GetSiteConfigById(ctx context.Context, id string) (*site_config.SiteConfig, error)
}
