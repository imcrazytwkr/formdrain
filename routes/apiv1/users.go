package apiv1

import (
	"context"
	"net/http"

	"github.com/imcrazytwkr/formdrain/middleware"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/oapi-codegen/runtime/types"
)

var getCurrentUserUnauthorizedResponse = api.GetCurrentUser401JSONResponse{
	Status:  http.StatusUnauthorized,
	Message: http.StatusText(http.StatusUnauthorized),
}

func (r *apiV1Router) GetCurrentUser(ctx context.Context, _ api.GetCurrentUserRequestObject) (api.GetCurrentUserResponseObject, error) {
	session, ok := middleware.SessionFromContext(ctx)
	if !ok {
		return &getCurrentUserUnauthorizedResponse, nil
	}

	account, err := r.accounts.GetByID(ctx, session.AccountID)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return &getCurrentUserUnauthorizedResponse, nil
	}

	return api.GetCurrentUser200JSONResponse{
		Id:    account.ID,
		Email: types.Email(account.Email),
	}, nil
}
