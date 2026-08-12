package apiv1

import (
	"context"
	"net/http"
	"time"

	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/mappers"
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

	var afterCreatedAt time.Time
	var afterID string
	if req.Params.Cursor != nil && len(*req.Params.Cursor) > 0 {
		afterCreatedAt, afterID, err = decodeCreatedAtCursor(*req.Params.Cursor)
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

	items := make([]api.FormResponse, len(rows))
	for i, row := range rows {
		items[i], err = mappers.FormResponse(row)
		if err != nil {
			return nil, err
		}
	}

	res := api.ListFormResponses200JSONResponse{Items: items}
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		next := encodeCreatedAtCursor(last.CreatedAt, last.Id)
		res.NextCursor = &next
	}

	return res, nil
}
