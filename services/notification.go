package services

import (
	"context"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

type NotificationService interface {
	Send(ctx context.Context, config fc.NotifiersConfig, form map[string]any) error
	SendAsync(ctx context.Context, config fc.NotifiersConfig, form map[string]any)
}
