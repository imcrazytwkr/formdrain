package form_config

import (
	"github.com/imcrazytwkr/formdrain/models/form_config/brevo"
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
)

type NotifiersConfig struct {
	Discord *discord.DiscordConfig `json:"discord"`
	Brevo   *brevo.BrevoConfig     `json:"brevo"`
}
