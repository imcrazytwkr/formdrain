package notifiers

import "github.com/imcrazytwkr/formdrain/models/form_config/brevo"

type BrevoNotifier interface {
	Send(config *brevo.BrevoConfig, form map[string]any) error
}
