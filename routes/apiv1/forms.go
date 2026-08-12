package apiv1

import (
	"context"
	"net/http"

	"github.com/imcrazytwkr/formdrain/middleware"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/mappers"
)

const defaultFormListLimit = 50
const maxFormListLimit = 100

var listFormsUnauthorized = api.ListForms401JSONResponse{
	Status:  http.StatusUnauthorized,
	Message: http.StatusText(http.StatusUnauthorized),
}

var listFormsBadRequest = api.ListForms400JSONResponse{
	Status:  http.StatusBadRequest,
	Message: http.StatusText(http.StatusBadRequest),
}

var getFormConfigNotFound = api.GetFormConfig404JSONResponse{
	Status:  http.StatusNotFound,
	Message: http.StatusText(http.StatusNotFound),
}

func (r *apiV1Router) GetFormConfig(ctx context.Context, req api.GetFormConfigRequestObject) (api.GetFormConfigResponseObject, error) {
	sess, ok := middleware.SessionFromContext(ctx)
	if !ok {
		return getFormConfigNotFound, nil
	}

	formConfig, err := r.forms.GetFormConfigById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if formConfig == nil {
		return getFormConfigNotFound, nil
	}

	siteConfig, err := r.sites.GetSiteConfigById(ctx, formConfig.SiteId)
	if err != nil {
		return nil, err
	}

	if siteConfig == nil || siteConfig.OwnerId != sess.AccountID {
		return getFormConfigNotFound, nil
	}

	res, err := mappers.FormConfig(formConfig)
	if err != nil {
		return nil, err
	}

	return api.GetFormConfig200JSONResponse(res), nil
}

func (r *apiV1Router) ListForms(ctx context.Context, req api.ListFormsRequestObject) (api.ListFormsResponseObject, error) {
	session, ok := middleware.SessionFromContext(ctx)
	if !ok {
		return listFormsUnauthorized, nil
	}

	sort := api.ListFormsParamsSortId
	if req.Params.Sort != nil {
		if !req.Params.Sort.Valid() {
			return listFormsBadRequest, nil
		}

		sort = *req.Params.Sort
	}

	limit := defaultFormListLimit
	if req.Params.Limit != nil {
		limit = min(*req.Params.Limit, maxFormListLimit)
		if limit < 1 {
			limit = defaultFormListLimit
		}
	}

	fetchLimit := limit + 1
	var rows []*fc.FormListItem
	var err error

	var siteId int64
	if req.Params.SiteId != nil {
		siteId = *req.Params.SiteId
	}

	if siteId > 0 {
		// Faster path
		site, err := r.sites.GetSiteConfigById(ctx, siteId)
		if err != nil {
			return nil, err
		}

		if site == nil || site.OwnerId != session.AccountID {
			return api.ListForms200JSONResponse{Items: nil}, nil
		}

		var afterID int64

		var cursor string
		if req.Params.Cursor != nil {
			cursor = *req.Params.Cursor
		}

		if len(cursor) > 0 {
			afterID, err = decodeIDCursor(cursor)
			if err != nil {
				return listFormsBadRequest, nil
			}
		}

		rows, err = r.forms.ListFormsBySiteID(ctx, *req.Params.SiteId, afterID, fetchLimit)
	} else {
		var afterID int64
		var afterSiteID int64
		var afterHostname string

		var cursor string
		if req.Params.Cursor != nil {
			cursor = *req.Params.Cursor
		}

		if len(cursor) > 0 {
			switch sort {
			case api.ListFormsParamsSortSiteId:
				afterSiteID, afterID, err = decodeSiteIDCursor(cursor)
			case api.ListFormsParamsSortHostname:
				afterHostname, afterID, err = decodeHostnameCursor(cursor)
			case api.ListFormsParamsSortId:
				afterID, err = decodeIDCursor(cursor)
			default:
				return listFormsBadRequest, nil
			}

			if err != nil {
				return listFormsBadRequest, nil
			}
		}

		rows, err = r.forms.ListFormsByOwnerID(ctx, session.AccountID, string(sort), afterID, afterSiteID, afterHostname, fetchLimit)
	}

	if err != nil {
		return nil, err
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	items := make([]api.FormListItem, len(rows))
	for i, row := range rows {
		items[i] = api.FormListItem{
			Id:     row.Id,
			SiteId: row.SiteId,
		}
	}

	res := api.ListForms200JSONResponse{Items: items}
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		var next string
		switch sort {
		case api.ListFormsParamsSortSiteId:
			next = encodeSiteIDCursor(last.SiteId, last.Id)
		case api.ListFormsParamsSortHostname:
			next = encodeHostnameCursor(last.Hostname, last.Id)
		case api.ListFormsParamsSortId:
			fallthrough
		default:
			next = encodeIDCursor(last.Id)
		}
		res.NextCursor = &next
	}

	return res, nil
}
