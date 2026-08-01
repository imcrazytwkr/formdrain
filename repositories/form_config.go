package repositories

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/form_config"
)

type FormConfigRepository interface {
	GetFormConfigById(ctx context.Context, id int64) (*form_config.FormConfig, error)
}
