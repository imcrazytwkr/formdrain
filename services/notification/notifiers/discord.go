package notifiers

import (
	"context"

	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
)

type DiscordNotifier interface {
	Send(ctx context.Context, config *discord.DiscordConfig, form map[string]any) error
}
