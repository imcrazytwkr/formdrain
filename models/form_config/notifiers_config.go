package form_config

import (
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
	"github.com/imcrazytwkr/formdrain/models/form_config/sendinblue"
)

type NotifiersConfig struct {
	Discord    *discord.DiscordConfig       `bson:"discord"`
	Sendinblue *sendinblue.SendinblueConfig `bson:"sendinblue"`
}
