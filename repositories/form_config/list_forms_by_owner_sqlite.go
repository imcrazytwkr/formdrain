package form_config

import (
	"context"
	"database/sql"
	"strings"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

const listFormsByOwnerBase = `
SELECT f.id, f.site_id, s.hostname
FROM sites s
	INNER JOIN forms f ON f.site_id = s.id
WHERE s.owner_id = ?
`

func (r *sqliteFormConfigRepository) ListFormsByOwnerID(ctx context.Context, ownerID int64, sort string, afterID int64, afterSiteID int64, afterHostname string, limit int) ([]*fc.FormListItem, error) {
	if ownerID < 1 || limit < 1 {
		return []*fc.FormListItem{}, nil
	}

	var qb strings.Builder
	qb.WriteString(listFormsByOwnerBase)

	var args []any
	args = append(args, ownerID)

	switch sort {
	case "site_id":
		if afterSiteID > 0 || afterID > 0 {
			qb.WriteString(" AND (f.site_id, f.id) > (?, ?)")
			args = append(args, afterSiteID, afterID)
		}

		qb.WriteString(" ORDER BY f.site_id ASC, f.id ASC")
	case "hostname":
		if len(afterHostname) > 0 || afterID > 0 {
			qb.WriteString(" AND (s.hostname, f.id) > (?, ?)")
			args = append(args, afterHostname, afterID)
		}

		qb.WriteString(" ORDER BY s.hostname ASC, f.id ASC")
	case "id":
		fallthrough
	default:
		if afterID > 0 {
			qb.WriteString(" AND f.id > ?")
			args = append(args, afterID)
		}

		qb.WriteString(" ORDER BY f.id")
	}

	qb.WriteString(" LIMIT ?")
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, qb.String(), args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanFormListItemsWithHostname(rows, limit)
}

func scanFormListItemsWithHostname(rows *sql.Rows, capacity int) (out []*fc.FormListItem, err error) {
	out = make([]*fc.FormListItem, 0, capacity)
	for rows.Next() {
		var item fc.FormListItem
		err = rows.Scan(&item.Id, &item.SiteId, &item.Hostname)
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
