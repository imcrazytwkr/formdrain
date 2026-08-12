package apiv1

import (
	"context"
	"net/http"

	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/models/site_config"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/cursors"
)

const defaultSiteListLimit = 50
const maxSiteListLimit = 100

var listSitesUnauthorized = api.ListSites401JSONResponse{
	Status:  http.StatusUnauthorized,
	Message: http.StatusText(http.StatusUnauthorized),
}

var listSitesBadRequest = api.ListSites400JSONResponse{
	Status:  http.StatusBadRequest,
	Message: http.StatusText(http.StatusBadRequest),
}

var getSiteConfigNotFound = api.GetSiteConfig404JSONResponse{
	Status:  http.StatusNotFound,
	Message: http.StatusText(http.StatusNotFound),
}

func (r *apiV1Router) GetSiteConfig(ctx context.Context, req api.GetSiteConfigRequestObject) (api.GetSiteConfigResponseObject, error) {
	sess, ok := middleware.SessionFromContext(ctx)
	if !ok {
		return getSiteConfigNotFound, nil
	}

	config, err := r.sites.GetSiteConfigById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if config == nil || config.OwnerId != sess.AccountID {
		return getSiteConfigNotFound, nil
	}

	return api.GetSiteConfig200JSONResponse{
		Id:       config.SiteId,
		Hostname: config.Hostname,
		OwnerId:  config.OwnerId,
	}, nil
}

func (r *apiV1Router) ListSites(ctx context.Context, req api.ListSitesRequestObject) (api.ListSitesResponseObject, error) {
	sess, ok := middleware.SessionFromContext(ctx)
	if !ok {
		return listSitesUnauthorized, nil
	}

	sort := api.ListSitesParamsSortId
	if req.Params.Sort != nil {
		if !req.Params.Sort.Valid() {
			return listSitesBadRequest, nil
		}

		sort = *req.Params.Sort
	}

	limit := defaultSiteListLimit
	if req.Params.Limit != nil {
		limit = min(*req.Params.Limit, maxSiteListLimit)
		if limit < 1 {
			limit = defaultSiteListLimit
		}
	}

	fetchLimit := limit + 1
	var rows []*site_config.SiteConfig
	var err error

	switch sort {
	case api.ListSitesParamsSortId:
		var afterID int64
		if req.Params.Cursor != nil && len(*req.Params.Cursor) > 0 {
			afterID, err = cursors.DecodeIDCursor(*req.Params.Cursor)
			if err != nil {
				return listSitesBadRequest, nil
			}
		}
		rows, err = r.sites.ListByOwnerIDAfterID(ctx, sess.AccountID, afterID, fetchLimit)
	case api.ListSitesParamsSortHostname:
		var afterHostname string
		var afterID int64
		if req.Params.Cursor != nil && len(*req.Params.Cursor) > 0 {
			afterHostname, afterID, err = cursors.DecodeHostnameCursor(*req.Params.Cursor)
			if err != nil {
				return listSitesBadRequest, nil
			}
		}
		rows, err = r.sites.ListByOwnerIDAfterHostname(ctx, sess.AccountID, afterHostname, afterID, fetchLimit)
	default:
		return listSitesBadRequest, nil
	}
	if err != nil {
		return nil, err
	}

	hasMore := len(rows) > limit
	if hasMore {
		rows = rows[:limit]
	}

	items := make([]api.Site, len(rows))
	for i, row := range rows {
		items[i] = api.Site{
			Id:       row.SiteId,
			Hostname: row.Hostname,
			OwnerId:  row.OwnerId,
		}
	}

	resp := api.ListSites200JSONResponse{Items: items}
	if hasMore && len(rows) > 0 {
		last := rows[len(rows)-1]
		var next string
		switch sort {
		case api.ListSitesParamsSortId:
			next = cursors.EncodeIDCursor(last.SiteId)
		case api.ListSitesParamsSortHostname:
			next = cursors.EncodeHostnameCursor(last.Hostname, last.SiteId)
		}
		resp.NextCursor = &next
	}

	return resp, nil
}
