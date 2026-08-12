package form_response

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/netip"
	"strings"
	"time"

	fr "github.com/imcrazytwkr/formdrain/models/form_response"
)

const listFormResponsesByFormIDBase = `
SELECT id, form_id, created_at, schema_version, client_ip, payload
FROM form_responses
WHERE form_id = ?
`

func (r *sqliteFormResponseRepository) ListFormResponsesByFormID(ctx context.Context, formID int64, afterCreatedAt time.Time, afterID string, limit int) ([]*fr.FormResponse, error) {
	if formID < 1 || limit < 1 {
		return []*fr.FormResponse{}, nil
	}

	var qb strings.Builder
	qb.WriteString(listFormResponsesByFormIDBase)

	args := []any{formID}
	if len(afterID) > 0 {
		qb.WriteString(" AND (created_at < ? OR (created_at = ? AND id > ?))")
		createdAt := afterCreatedAt.UTC().Format(time.RFC3339)
		args = append(args, createdAt, createdAt, afterID)
	}

	qb.WriteString(" ORDER BY created_at DESC, id ASC LIMIT ?")
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, qb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanFormResponses(rows, limit)
}

func scanFormResponses(rows *sql.Rows, capacity int) ([]*fr.FormResponse, error) {
	out := make([]*fr.FormResponse, 0, capacity)
	for rows.Next() {
		var item fr.FormResponse
		var createdAt string
		var clientIP sql.NullString
		var payloadJSON string

		err := rows.Scan(
			&item.Id,
			&item.FormId,
			&createdAt,
			&item.SchemaVersion,
			&clientIP,
			&payloadJSON,
		)
		if err != nil {
			return nil, err
		}

		item.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}
		item.CreatedAt = item.CreatedAt.UTC()

		if clientIP.Valid {
			item.ClientIP, err = netip.ParseAddr(clientIP.String)
			if err != nil {
				return nil, err
			}
		}

		err = json.Unmarshal([]byte(payloadJSON), &item.Payload)
		if err != nil {
			return nil, err
		}
		if item.Payload == nil {
			item.Payload = map[string]any{}
		}

		out = append(out, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
