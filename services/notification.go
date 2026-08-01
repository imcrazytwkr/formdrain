package services

import fc "github.com/imcrazytwkr/formdrain/models/form_config"

type NotificationService interface {
	Send(config fc.NotifiersConfig, form map[string]any) error
}
