package repositories

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/form_response"
)

type FormResponseRepository interface {
	SaveFormResponse(ctx context.Context, response *form_response.FormResponse) (string, error)
}
