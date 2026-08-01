package form

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/models/site_config"
)

type formConfigCtxKey struct{}
type siteConfigCtxKey struct{}

func withFormConfig(ctx context.Context, cfg *form_config.FormConfig) context.Context {
	return context.WithValue(ctx, formConfigCtxKey{}, cfg)
}

func withSiteConfig(ctx context.Context, cfg *site_config.SiteConfig) context.Context {
	return context.WithValue(ctx, siteConfigCtxKey{}, cfg)
}

func getFormConfig(ctx context.Context) (*form_config.FormConfig, bool) {
	cfg, ok := ctx.Value(formConfigCtxKey{}).(*form_config.FormConfig)
	return cfg, ok
}

func getSiteConfig(ctx context.Context) (*site_config.SiteConfig, bool) {
	cfg, ok := ctx.Value(siteConfigCtxKey{}).(*site_config.SiteConfig)
	return cfg, ok
}
