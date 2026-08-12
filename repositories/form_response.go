package repositories

import (
	"context"
	"time"

	"github.com/imcrazytwkr/formdrain/models/form_response"
)

type FormResponseRepository interface {
	SaveFormResponse(ctx context.Context, response *form_response.FormResponse) (string, error)
	ListFormResponsesByFormID(ctx context.Context, formID int64, afterCreatedAt time.Time, afterID string, limit int) ([]*form_response.FormResponse, error)
}
