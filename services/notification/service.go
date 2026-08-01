package notification

import (
	"errors"
	"net/http"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/services"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers"
	dn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/discord"
	sn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/sendinblue"
)

type httpNotificationService struct {
	discordNotifier    notifiers.DiscordNotifier
	sendinblueNotifier notifiers.SendinblueNotifier
}

func NewHttpNotificationService(httpClient *http.Client) services.NotificationService {
	return &httpNotificationService{
		discordNotifier:    dn.NewDiscordNotifier("discord_username", "discord_avatar", httpClient),
		sendinblueNotifier: sn.NewSendinblueNotifier("sender_name", "sender_email", httpClient),
	}
}

func (s *httpNotificationService) Send(config *fc.NotifiersConfig, form map[string]any) error {
	var errs []error

	if config.Discord != nil {
		err := s.discordNotifier.Send(config.Discord, form)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if config.Sendinblue != nil {
		err := s.sendinblueNotifier.Send(config.Sendinblue, form)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
