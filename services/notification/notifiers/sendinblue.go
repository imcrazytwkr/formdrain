package notifiers

import "github.com/imcrazytwkr/formdrain/models/form_config/sendinblue"

type SendinblueNotifier interface {
	Send(config *sendinblue.SendinblueConfig, form map[string]any) error
}
