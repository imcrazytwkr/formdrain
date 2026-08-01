package discord

import (
	"github.com/imcrazytwkr/formdrain/models/common"
)

type DiscordConfig struct {
	Webhooks []*WebhookKey    `json:"webhooks"`
	Author   *Author          `json:"author"`
	Title    string           `json:"title,omitempty"`
	Url      string           `json:"url,omitempty"`
	Template *common.Template `json:"template"`
	Color    int              `json:"color,omitempty"`
}

func (c *DiscordConfig) RenderContent(form map[string]any) string {
	return c.Template.ExecuteString(form)
}
