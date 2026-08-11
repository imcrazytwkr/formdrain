package repositories

import (
	"context"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

type FormConfigRepository interface {
	GetFormConfigById(ctx context.Context, id int64) (*fc.FormConfig, error)
	// ListFormsByOwnerID lists all forms across all sites owned by the user.
	// sort must be one of: "id", "site_id", "hostname".
	ListFormsByOwnerID(ctx context.Context, ownerID int64, sort string, afterID int64, afterSiteID int64, afterHostname string, limit int) ([]*fc.FormListItem, error)
	// ListFormsBySiteID lists forms for a specific site. Ownership is NOT checked here —
	// the caller must verify the site belongs to the user before calling.
	// sort is ignored when site_id is set (forms are always ordered by id).
	ListFormsBySiteID(ctx context.Context, siteID int64, afterID int64, limit int) ([]*fc.FormListItem, error)
}
