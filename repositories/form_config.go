package repositories

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/form_config"
)

type FormConfigRepository interface {
	GetFormConfigById(ctx context.Context, id string) (*form_config.FormConfig, error)
}
