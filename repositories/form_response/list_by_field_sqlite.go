package form_response

import (
	"context"
	"strings"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	fr "github.com/imcrazytwkr/formdrain/models/form_response"
	"github.com/imcrazytwkr/formdrain/utils/collate"
)

// I would prefer to use json-extract rather than hack in a join here. However,
// assembling a safe jsonpath is more risky than what performance benefirs are
// worth
const listFormResponsesByFieldBase = `
SELECT fr.id, fr.form_id, fr.created_at, fr.schema_version, fr.client_ip, fr.payload
FROM form_responses fr
LEFT JOIN json_each(fr.payload) je ON je.key = ?
WHERE fr.form_id = ?
`

/**
 * The outright BS approach with repeating CASE-WHEN is necessary only because
 * SQLite is not ANSI-SQL complaint when it comes to field evaluation and can
 * only filter by actual commands or inlined SQL expressions only using these
 * columns.
 */
func (r *sqliteFormResponseRepository) ListFormResponsesByField(ctx context.Context, formID int64, field string, fieldType fc.FieldType, desc bool, afterNull bool, afterValue any, afterID string, limit int) ([]*fr.FormResponse, error) {
	if formID < 1 || limit < 1 || len(field) < 1 {
		return []*fr.FormResponse{}, nil
	}

	expr, ok := sortValueExpr(fieldType)
	if !ok {
		return []*fr.FormResponse{}, nil
	}

	var qb strings.Builder
	qb.WriteString(listFormResponsesByFieldBase)
	args := []any{field, formID}

	if len(afterID) > 0 {
		if afterNull {
			qb.WriteString(" AND ")
			qb.WriteString(expr)
			qb.WriteString(" IS NULL AND fr.id > ?")
			args = append(args, afterID)
		} else {
			bound, ok := bindSortValue(fieldType, afterValue)
			if !ok {
				return []*fr.FormResponse{}, nil
			}

			qb.WriteString(" AND (")
			qb.WriteString(expr)
			qb.WriteString(" IS NULL OR ")
			qb.WriteString(expr)
			qb.WriteByte(' ')
			if desc {
				qb.WriteByte('<')
			} else {
				qb.WriteByte('>')
			}
			qb.WriteString(" ? OR (")
			qb.WriteString(expr)
			qb.WriteString(" = ? AND fr.id > ?))")
			args = append(args, bound, bound, afterID)
		}
	}

	qb.WriteString(" ORDER BY (")
	qb.WriteString(expr)
	qb.WriteString(" IS NULL) ASC, ")
	qb.WriteString(expr)
	if desc {
		qb.WriteString(" DESC")
	} else {
		qb.WriteString(" ASC")
	}
	qb.WriteString(", fr.id ASC LIMIT ?")
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, qb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanFormResponses(rows, limit)
}

func sortValueExpr(fieldType fc.FieldType) (string, bool) {
	switch fieldType {
	case fc.FieldTypeString:
		return "CASE WHEN je.type = 'text' THEN je.atom ELSE NULL END", true
	case fc.FieldTypeNumber:
		return "CASE WHEN je.type IN ('integer', 'real') THEN je.atom ELSE NULL END", true
	case fc.FieldTypeBoolean:
		return "CASE WHEN je.type IN ('true', 'false') THEN je.atom ELSE NULL END", true
	default:
		return "", false
	}
}

func bindSortValue(fieldType fc.FieldType, value any) (any, bool) {
	switch fieldType {
	case fc.FieldTypeString:
		s, ok := value.(string)
		return s, ok
	case fc.FieldTypeNumber:
		return collate.NumberToFloat64(value)
	case fc.FieldTypeBoolean:
		return collate.BoolToInt(value)
	default:
		return nil, false
	}
}
