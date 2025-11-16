package discord

import (
	"github.com/imcrazytwkr/formdrain/models/common"
)

type DiscordConfig struct {
	Webhooks []*WebhookKey    `bson:"webhooks"`
	Author   *Author          `bson:"author"`
	Title    string           `bson:"title,omitempty"`
	Url      string           `bson:"url,omitempty"`
	Template *common.Template `bson:"template"`
	Color    int              `bson:"color,omitempty"`
}

func (c *DiscordConfig) RenderContent(form map[string]any) string {
	return c.Template.ExecuteString(form)
}
