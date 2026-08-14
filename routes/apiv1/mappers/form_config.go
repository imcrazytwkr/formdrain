package mappers

import (
	"github.com/imcrazytwkr/formdrain/models/common"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/models/form_config/brevo"
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/oapi-codegen/runtime/types"
)

func FormConfig(config *fc.FormConfig) (api.FormConfig, error) {
	return api.FormConfig{
		Id:            config.FormId,
		SiteId:        config.SiteId,
		CaptchaType:   api.CaptchaType(config.CaptchaType.String()),
		CaptchaField:  optString(config.CaptchaField),
		RedirectTo:    optString(config.RedirectTo),
		FieldSchema:   fieldSchema(config.FieldSchema),
		SchemaVersion: config.SchemaVersion,
		Notifiers:     notifiersConfig(config.Notifiers),
	}, nil
}

func fieldSchema(src fc.FieldSchema) api.FieldSchema {
	fields := make([]api.Field, len(src.Fields))
	for i, f := range src.Fields {
		fields[i] = schemaField(f)
	}

	return api.FieldSchema{
		Fields: fields,
	}
}

func schemaField(src fc.Field) api.Field {
	return api.Field{
		Name:     src.Name,
		Type:     api.FieldType(src.Type),
		Required: src.Required,
		Items:    schemaFieldItems(src.Items),
	}
}

func schemaFieldItems(src *fc.FieldItems) *api.FieldItems {
	if src == nil {
		return nil
	}

	return &api.FieldItems{Type: api.FieldType(src.Type)}
}

func notifiersConfig(src fc.NotifiersConfig) api.NotifiersConfig {
	return api.NotifiersConfig{
		Brevo:   brevoConfig(src.Brevo),
		Discord: discordConfig(src.Discord),
	}
}

func brevoConfig(src *brevo.BrevoConfig) *api.BrevoConfig {
	if src == nil {
		return nil
	}

	recipients := make([]api.EmailContact, len(src.Recipients))
	for i, r := range src.Recipients {
		recipients[i] = brevoEmailContact(r)
	}

	return &api.BrevoConfig{
		Subject:    src.Subject,
		Template:   templateString(src.Template),
		Recipients: recipients,
	}
}

func brevoEmailContact(src *brevo.EmailContact) api.EmailContact {
	if src == nil {
		return api.EmailContact{}
	}

	return api.EmailContact{
		Email: types.Email(src.Address),
		Name:  optString(src.Name),
	}
}

func discordConfig(src *discord.DiscordConfig) *api.DiscordConfig {
	if src == nil {
		return nil
	}

	return &api.DiscordConfig{
		Title:    optString(src.Title),
		Url:      optString(src.Url),
		Color:    optInt(src.Color),
		Template: templateString(src.Template),
		Webhooks: discordWebhooks(src.Webhooks),
		Author:   discordAuthor(src.Author),
	}
}

func discordWebhooks(src []*discord.WebhookKey) *[]api.WebhookKey {
	if src == nil {
		return nil
	}

	webhooks := make([]api.WebhookKey, len(src))
	for i, wh := range src {
		if wh == nil {
			continue
		}
		webhooks[i] = api.WebhookKey{
			Snowflake: wh.Snowflake,
			Token:     wh.Token,
		}
	}

	return &webhooks
}

func discordAuthor(src *discord.Author) *api.Author {
	if src == nil {
		return nil
	}

	return &api.Author{
		Name:    optString(src.Name),
		Url:     optString(src.Url),
		IconUrl: optString(src.Icon),
	}
}

func templateString(t *common.Template) *string {
	if t == nil {
		return nil
	}

	s := t.String()
	return &s
}

func optString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optInt(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}
