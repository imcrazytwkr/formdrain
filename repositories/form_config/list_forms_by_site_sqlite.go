package form_config

import (
	"context"
	"database/sql"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

const listFormsBySiteID = `
SELECT id, site_id
FROM forms
WHERE site_id = ? AND id > ?
ORDER BY id ASC LIMIT ?
`

// @NOTE: extra owner check is not necessary because we have already checkede it when fetching SiteConfig
func (r *sqliteFormConfigRepository) ListFormsBySiteID(ctx context.Context, siteID int64, afterID int64, limit int) ([]*fc.FormListItem, error) {
	if siteID < 1 || limit < 1 {
		return nil, nil
	}

	rows, err := r.db.QueryContext(ctx, listFormsBySiteID, siteID, afterID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanFormListItems(rows, limit)
}

func scanFormListItems(rows *sql.Rows, capacity int) (out []*fc.FormListItem, err error) {
	out = make([]*fc.FormListItem, 0, capacity)
	for rows.Next() {
		var item fc.FormListItem
		err = rows.Scan(&item.Id, &item.SiteId)
		if err != nil {
			return nil, err
		}

		out = append(out, &item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return out, nil
}
