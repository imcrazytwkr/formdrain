package notifiers

import "github.com/imcrazytwkr/formdrain/models/form_config/discord"

type DiscordNotifier interface {
	Send(config *discord.DiscordConfig, form map[string]any) error
}
