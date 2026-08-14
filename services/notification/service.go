package notification

import (
	"context"
	"errors"
	"net/http"

	"github.com/imcrazytwkr/formdrain/models/config"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/services"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers"
	bn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/brevo"
	dn "github.com/imcrazytwkr/formdrain/services/notification/notifiers/discord"
	"github.com/imcrazytwkr/formdrain/utils/logutil"
	"github.com/rs/zerolog"
)

type httpNotificationService struct {
	discordNotifier notifiers.DiscordNotifier
	brevoNotifier   notifiers.BrevoNotifier
}

func NewHttpNotificationService(
	httpClient *http.Client,
	cfg config.NotifiersConfig,
	brevoAPIKey string,
) services.NotificationService {
	return &httpNotificationService{
		discordNotifier: dn.NewDiscordNotifier(cfg.Discord, httpClient),
		brevoNotifier:   bn.NewBrevoNotifier(cfg.Brevo, brevoAPIKey, httpClient),
	}
}

func (s *httpNotificationService) Send(ctx context.Context, config fc.NotifiersConfig, form map[string]any) error {
	var errs []error

	if config.Discord != nil {
		err := s.discordNotifier.Send(ctx, config.Discord, form)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if config.Brevo != nil {
		err := s.brevoNotifier.Send(ctx, config.Brevo, form)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (s *httpNotificationService) sendAsync(ctx context.Context, config fc.NotifiersConfig, form map[string]any) {
	notifyCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), notificationTimeout)
	defer cancel()

	err := s.Send(notifyCtx, config, form)
	if err != nil {
		logutil.UnwrapErr(zerolog.Ctx(notifyCtx).Error(), err).Msg("failed to send notifications")
	}
}

func (s *httpNotificationService) SendAsync(ctx context.Context, config fc.NotifiersConfig, form map[string]any) {
	go s.sendAsync(ctx, config, form)
}
