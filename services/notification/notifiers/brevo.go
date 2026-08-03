package notifiers

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/form_config/brevo"
)

type BrevoNotifier interface {
	Send(ctx context.Context, config *brevo.BrevoConfig, form map[string]any) error
}
