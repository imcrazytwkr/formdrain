package repositories

import "context"

type FormResponseRepository interface {
	SaveFormResponse(ctx context.Context, form map[string]any) (string, error)
}
