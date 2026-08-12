package apiv1

import (
	"context"
	"net/http"

	"github.com/imcrazytwkr/formdrain/middleware"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/cursors"
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
	sess, ok := middleware.SessionFromContext(ctx)
	if !ok {
		return listFormsUnauthorized, nil
	}

	limit := defaultFormListLimit
	if req.Params.Limit != nil {
		limit = min(*req.Params.Limit, maxFormListLimit)
		if limit < 1 {
			limit = defaultFormListLimit
		}
	}

	var cursor string
	if req.Params.Cursor != nil {
		cursor = *req.Params.Cursor
	}

	var siteId int64
	if req.Params.SiteId != nil {
		siteId = *req.Params.SiteId
	}

	if siteId > 0 {
		// Site filter skips the owner join; ownership is checked on the site first.
		return r.listFormsForSite(ctx, sess.AccountID, siteId, cursor, limit)
	}

	sort := api.ListFormsParamsSortId
	if req.Params.Sort != nil {
		if !req.Params.Sort.Valid() {
			return listFormsBadRequest, nil
		}

		sort = *req.Params.Sort
	}

	return r.listFormsForOwner(ctx, sess.AccountID, sort, cursor, limit)
}

func (r *apiV1Router) listFormsForSite(ctx context.Context, accountID, siteID int64, cursor string, limit int) (api.ListFormsResponseObject, error) {
	site, err := r.sites.GetSiteConfigById(ctx, siteID)
	if err != nil {
		return nil, err
	}

	if site == nil || site.OwnerId != accountID {
		return api.ListForms200JSONResponse{Items: nil}, nil
	}

	var afterID int64
	if len(cursor) > 0 {
		afterID, err = cursors.DecodeIDCursor(cursor)
		if err != nil {
			return listFormsBadRequest, nil
		}
	}

	rows, err := r.forms.ListFormsBySiteID(ctx, siteID, afterID, limit+1)
	if err != nil {
		return nil, err
	}

	return formListPage(rows, limit, api.ListFormsParamsSortId), nil
}

func (r *apiV1Router) listFormsForOwner(ctx context.Context, accountID int64, sort api.ListFormsParamsSort, cursor string, limit int) (api.ListFormsResponseObject, error) {
	var afterID int64
	var afterSiteID int64
	var afterHostname string
	var err error

	if len(cursor) > 0 {
		switch sort {
		case api.ListFormsParamsSortSiteId:
			afterSiteID, afterID, err = cursors.DecodeSiteIDCursor(cursor)
		case api.ListFormsParamsSortHostname:
			afterHostname, afterID, err = cursors.DecodeHostnameCursor(cursor)
		case api.ListFormsParamsSortId:
			afterID, err = cursors.DecodeIDCursor(cursor)
		default:
			return listFormsBadRequest, nil
		}

		if err != nil {
			return listFormsBadRequest, nil
		}
	}

	rows, err := r.forms.ListFormsByOwnerID(ctx, accountID, string(sort), afterID, afterSiteID, afterHostname, limit+1)
	if err != nil {
		return nil, err
	}

	return formListPage(rows, limit, sort), nil
}

func formListPage(rows []*fc.FormListItem, limit int, sort api.ListFormsParamsSort) api.ListForms200JSONResponse {
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
	if !hasMore || len(rows) < 1 {
		return res
	}

	last := rows[len(rows)-1]
	next := encodeFormListCursor(sort, last)
	res.NextCursor = &next
	return res
}

func encodeFormListCursor(sort api.ListFormsParamsSort, last *fc.FormListItem) string {
	switch sort {
	case api.ListFormsParamsSortSiteId:
		return cursors.EncodeSiteIDCursor(last.SiteId, last.Id)
	case api.ListFormsParamsSortHostname:
		return cursors.EncodeHostnameCursor(last.Hostname, last.Id)
	default:
		return cursors.EncodeIDCursor(last.Id)
	}
}
