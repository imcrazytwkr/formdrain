package notification

import (
	"errors"
	"net/http"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/services"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers"
	bn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/brevo"
	dn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/discord"
)

type httpNotificationService struct {
	discordNotifier notifiers.DiscordNotifier
	brevoNotifier   notifiers.BrevoNotifier
}

func NewHttpNotificationService(httpClient *http.Client) services.NotificationService {
	return &httpNotificationService{
		discordNotifier: dn.NewDiscordNotifier("discord_username", "discord_avatar", httpClient),
		brevoNotifier:   bn.NewBrevoNotifier("sender_name", "sender_email", "brevo_api_key", httpClient),
	}
}

func (s *httpNotificationService) Send(config fc.NotifiersConfig, form map[string]any) error {
	var errs []error

	if config.Discord != nil {
		err := s.discordNotifier.Send(config.Discord, form)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if config.Brevo != nil {
		err := s.brevoNotifier.Send(config.Brevo, form)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
