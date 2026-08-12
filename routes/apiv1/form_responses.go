package apiv1

import (
	"context"
	"net/http"
	"time"

	"github.com/imcrazytwkr/formdrain/middleware"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	fr "github.com/imcrazytwkr/formdrain/models/form_response"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/cursors"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/mappers"
	"github.com/imcrazytwkr/formdrain/utils/collate"
)

const defaultFormResponseListLimit = 50
const maxFormResponseListLimit = 100

var listFormResponsesUnauthorized = api.ListFormResponses401JSONResponse{
	Status:  http.StatusUnauthorized,
	Message: http.StatusText(http.StatusUnauthorized),
}

var listFormResponsesBadRequest = api.ListFormResponses400JSONResponse{
	Status:  http.StatusBadRequest,
	Message: http.StatusText(http.StatusBadRequest),
}

var listFormResponsesNotFound = api.ListFormResponses404JSONResponse{
	Status:  http.StatusNotFound,
	Message: http.StatusText(http.StatusNotFound),
}

func (r *apiV1Router) ListFormResponses(ctx context.Context, req api.ListFormResponsesRequestObject) (api.ListFormResponsesResponseObject, error) {
	sess, ok := middleware.SessionFromContext(ctx)
	if !ok {
		return listFormResponsesUnauthorized, nil
	}

	formConfig, err := r.forms.GetFormConfigById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if formConfig == nil {
		return listFormResponsesNotFound, nil
	}

	siteConfig, err := r.sites.GetSiteConfigById(ctx, formConfig.SiteId)
	if err != nil {
		return nil, err
	}

	if siteConfig == nil || siteConfig.OwnerId != sess.AccountID {
		return listFormResponsesNotFound, nil
	}

	limit := defaultFormResponseListLimit
	if req.Params.Limit != nil {
		limit = min(*req.Params.Limit, maxFormResponseListLimit)
		if limit < 1 {
			limit = defaultFormResponseListLimit
		}
	}

	var sortField string
	if req.Params.Sort != nil {
		sortField = *req.Params.Sort
	}

	if req.Params.Order != nil && !req.Params.Order.Valid() {
		return listFormResponsesBadRequest, nil
	}

	var cursor string
	if req.Params.Cursor != nil {
		cursor = *req.Params.Cursor
	}

	if len(sortField) > 0 {
		field, ok := sortableSchemaField(formConfig.FieldSchema, sortField)
		if !ok {
			return listFormResponsesBadRequest, nil
		}

		desc := false
		if req.Params.Order != nil {
			desc = *req.Params.Order == api.Desc
		}

		var afterNull bool
		var afterValue any
		var afterID string
		if len(cursor) > 0 {
			decoded, err := cursors.DecodeFieldCursor(cursor)
			if err != nil || decoded.Field != field.Name || decoded.Type != field.Type || decoded.Desc != desc {
				return listFormResponsesBadRequest, nil
			}

			afterNull = decoded.Null
			afterValue = decoded.Value
			afterID = decoded.ID
		}

		rows, err := r.responses.ListFormResponsesByField(ctx, req.Id, field.Name, field.Type, desc, afterNull, afterValue, afterID, limit+1)
		if err != nil {
			return nil, err
		}

		hasMore := len(rows) > limit
		if hasMore {
			rows = rows[:limit]
		}

		items, err := mapFormResponses(rows)
		if err != nil {
			return nil, err
		}

		res := api.ListFormResponses200JSONResponse{Items: items}
		if hasMore && len(rows) > 0 {
			last := rows[len(rows)-1]
			value, isNull := payloadSortValue(last.Payload, field.Name, field.Type)
			next, err := cursors.EncodeFieldCursor(cursors.FieldCursor{
				Field: field.Name,
				Type:  field.Type,
				Desc:  desc,
				Null:  isNull,
				Value: value,
				ID:    last.Id,
			})
			if err != nil {
				return nil, err
			}
			res.NextCursor = &next
		}
		return res, nil
	}

	var afterCreatedAt time.Time
	var afterID string
	if len(cursor) > 0 {
		afterCreatedAt, afterID, err = cursors.DecodeCreatedAtCursor(cursor)
		if err != nil {
			return listFormResponsesBadRequest, nil
		}
	}

	rows, err := r.responses.ListFormResponsesByFormID(ctx, req.Id, afterCreatedAt, afterID, limit+1)
	if err != nil {
		return nil, err
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	items, err := mapFormResponses(rows)
	if err != nil {
		return nil, err
	}

	res := api.ListFormResponses200JSONResponse{Items: items}
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		next := cursors.EncodeCreatedAtCursor(last.CreatedAt, last.Id)
		res.NextCursor = &next
	}
	return res, nil
}

func sortableSchemaField(schema fc.FieldSchema, name string) (fc.Field, bool) {
	for _, field := range schema.Fields {
		if field.Name != name {
			continue
		}
		switch field.Type {
		case fc.FieldTypeString, fc.FieldTypeNumber, fc.FieldTypeBoolean:
			return field, true
		default:
			return fc.Field{}, false
		}
	}
	return fc.Field{}, false
}

func mapFormResponses(rows []*fr.FormResponse) ([]api.FormResponse, error) {
	items := make([]api.FormResponse, len(rows))
	for i, row := range rows {
		item, err := mappers.FormResponse(row)
		if err != nil {
			return nil, err
		}
		items[i] = item
	}
	return items, nil
}

func payloadSortValue(payload map[string]any, field string, fieldType fc.FieldType) (any, bool) {
	if payload == nil {
		return nil, true
	}

	raw, ok := payload[field]
	if !ok || raw == nil {
		return nil, true
	}

	switch fieldType {
	case fc.FieldTypeString:
		s, ok := raw.(string)
		if !ok {
			return nil, true
		}

		return s, false
	case fc.FieldTypeNumber:
		f, ok := collate.NumberToFloat64(raw)
		if !ok {
			return nil, true
		}

		return f, false
	case fc.FieldTypeBoolean:
		b, ok := raw.(bool)
		if !ok {
			return nil, true
		}

		return b, false
	default:
		return nil, true
	}
}
