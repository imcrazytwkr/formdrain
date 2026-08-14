package mappers_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/models/form_config/brevo"
	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/api"
	"github.com/imcrazytwkr/formdrain/routes/apiv1/mappers"
	"github.com/oapi-codegen/runtime/types"
)

func TestFormConfig_Minimal(t *testing.T) {
	t.Parallel()

	src := &fc.FormConfig{
		FormId:        42,
		SiteId:        7,
		CaptchaType:   common.CaptchaTypeHcaptcha,
		SchemaVersion: 1,
		FieldSchema: fc.FieldSchema{
			Fields: nil,
		},
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	if got.Id != 42 {
		t.Fatalf("Id = %d, want 42", got.Id)
	}
	if got.SiteId != 7 {
		t.Fatalf("SiteId = %d, want 7", got.SiteId)
	}
	if got.CaptchaType != api.Hcaptcha {
		t.Fatalf("CaptchaType = %q, want %q", got.CaptchaType, api.Hcaptcha)
	}
	if got.SchemaVersion != 1 {
		t.Fatalf("SchemaVersion = %d, want 1", got.SchemaVersion)
	}
	if len(got.FieldSchema.Fields) != 0 {
		t.Fatalf("FieldSchema.Fields = %#v, want empty", got.FieldSchema.Fields)
	}
}

func TestFormConfig_CaptchaTypes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		src  common.CaptchaType
		want api.CaptchaType
	}{
		{"hcaptcha", common.CaptchaTypeHcaptcha, api.Hcaptcha},
		{"recaptcha", common.CaptchaTypeRecaptcha, api.Recaptcha},
		{"undefined", common.CaptchaTypeUndefined, api.CaptchaType("undefined")},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := mappers.FormConfig(&fc.FormConfig{CaptchaType: tc.src})
			if err != nil {
				t.Fatal(err)
			}
			if got.CaptchaType != tc.want {
				t.Fatalf("CaptchaType = %q, want %q", got.CaptchaType, tc.want)
			}
		})
	}
}

func TestFormConfig_OptionalScalars(t *testing.T) {
	t.Parallel()

	src := &fc.FormConfig{
		FormId:       1,
		CaptchaType:  common.CaptchaTypeRecaptcha,
		CaptchaField: "cf-turnstile-response",
		RedirectTo:   "https://example.com/thanks",
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	if got.CaptchaField == nil || *got.CaptchaField != "cf-turnstile-response" {
		t.Fatalf("CaptchaField = %#v", got.CaptchaField)
	}
	if got.RedirectTo == nil || *got.RedirectTo != "https://example.com/thanks" {
		t.Fatalf("RedirectTo = %#v", got.RedirectTo)
	}
}

func TestFormConfig_FieldSchema(t *testing.T) {
	t.Parallel()

	src := &fc.FormConfig{
		FormId: 1,
		FieldSchema: fc.FieldSchema{
			Fields: []fc.Field{
				{Name: "email", Type: fc.FieldTypeString, Required: true},
				{Name: "age", Type: fc.FieldTypeNumber, Required: false},
				{Name: "active", Type: fc.FieldTypeBoolean, Required: true},
				{
					Name:     "tags",
					Type:     fc.FieldTypeArray,
					Required: false,
					Items:    &fc.FieldItems{Type: fc.FieldTypeString},
				},
			},
		},
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(got.FieldSchema.Fields) != 4 {
		t.Fatalf("len(Fields) = %d, want 4", len(got.FieldSchema.Fields))
	}

	email := got.FieldSchema.Fields[0]
	if email.Name != "email" || email.Type != api.String || !email.Required {
		t.Fatalf("email field = %#v", email)
	}
	if email.Items != nil {
		t.Fatalf("email.Items = %#v, want nil", email.Items)
	}

	tags := got.FieldSchema.Fields[3]
	if tags.Name != "tags" || tags.Type != api.Array || tags.Required {
		t.Fatalf("tags field = %#v", tags)
	}
	if tags.Items == nil || tags.Items.Type != api.String {
		t.Fatalf("tags.Items = %#v", tags.Items)
	}
}

func TestFormConfig_BrevoNotifier(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("Hello {{name}}")
	if err != nil {
		t.Fatal(err)
	}

	src := &fc.FormConfig{
		FormId: 1,
		Notifiers: fc.NotifiersConfig{
			Brevo: &brevo.BrevoConfig{
				Subject:  "New submission",
				Template: tmpl,
				Recipients: []*brevo.EmailContact{
					{Name: "Alice", Address: "alice@example.com"},
					{Address: "bob@example.com"},
				},
			},
		},
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	if got.Notifiers.Brevo == nil {
		t.Fatal("expected Brevo config")
	}
	if got.Notifiers.Discord != nil {
		t.Fatalf("Discord = %#v, want nil", got.Notifiers.Discord)
	}

	b := got.Notifiers.Brevo
	if b.Subject != "New submission" {
		t.Fatalf("Subject = %q", b.Subject)
	}
	if b.Template == nil || *b.Template != "Hello {{name}}" {
		t.Fatalf("Template = %#v", b.Template)
	}
	if len(b.Recipients) != 2 {
		t.Fatalf("len(Recipients) = %d, want 2", len(b.Recipients))
	}

	if b.Recipients[0].Email != types.Email("alice@example.com") {
		t.Fatalf("Recipients[0].Email = %q", b.Recipients[0].Email)
	}
	if b.Recipients[0].Name == nil || *b.Recipients[0].Name != "Alice" {
		t.Fatalf("Recipients[0].Name = %#v", b.Recipients[0].Name)
	}

	if b.Recipients[1].Email != types.Email("bob@example.com") {
		t.Fatalf("Recipients[1].Email = %q", b.Recipients[1].Email)
	}
}

func TestFormConfig_BrevoNilTemplate(t *testing.T) {
	t.Parallel()

	src := &fc.FormConfig{
		FormId: 1,
		Notifiers: fc.NotifiersConfig{
			Brevo: &brevo.BrevoConfig{
				Subject:    "No body",
				Recipients: []*brevo.EmailContact{},
			},
		},
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	if got.Notifiers.Brevo == nil {
		t.Fatal("expected Brevo config")
	}
	if got.Notifiers.Brevo.Template != nil {
		t.Fatalf("Template = %#v, want nil", got.Notifiers.Brevo.Template)
	}
}

func TestFormConfig_DiscordNotifier(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("msg={{email}}")
	if err != nil {
		t.Fatal(err)
	}

	src := &fc.FormConfig{
		FormId: 1,
		Notifiers: fc.NotifiersConfig{
			Discord: &discord.DiscordConfig{
				Title:    "Form alert",
				Url:      "https://example.com/form",
				Color:    0x112233,
				Template: tmpl,
				Webhooks: []*discord.WebhookKey{
					{Snowflake: "1234567890", Token: "wh-token"},
				},
				Author: &discord.Author{
					Name: "Bot",
					Url:  "https://example.com/bot",
					Icon: "https://example.com/icon.png",
				},
			},
		},
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	if got.Notifiers.Discord == nil {
		t.Fatal("expected Discord config")
	}
	if got.Notifiers.Brevo != nil {
		t.Fatalf("Brevo = %#v, want nil", got.Notifiers.Brevo)
	}

	d := got.Notifiers.Discord
	if d.Title == nil || *d.Title != "Form alert" {
		t.Fatalf("Title = %#v", d.Title)
	}
	if d.Url == nil || *d.Url != "https://example.com/form" {
		t.Fatalf("Url = %#v", d.Url)
	}
	if d.Color == nil || *d.Color != 0x112233 {
		t.Fatalf("Color = %#v", d.Color)
	}
	if d.Template == nil || *d.Template != "msg={{email}}" {
		t.Fatalf("Template = %#v", d.Template)
	}

	if d.Webhooks == nil || len(*d.Webhooks) != 1 {
		t.Fatalf("Webhooks = %#v", d.Webhooks)
	}
	wh := (*d.Webhooks)[0]
	if wh.Snowflake != "1234567890" || wh.Token != "wh-token" {
		t.Fatalf("webhook = %#v", wh)
	}

	if d.Author == nil {
		t.Fatal("expected Author")
	}
	if d.Author.Name == nil || *d.Author.Name != "Bot" {
		t.Fatalf("Author.Name = %#v", d.Author.Name)
	}
	if d.Author.Url == nil || *d.Author.Url != "https://example.com/bot" {
		t.Fatalf("Author.Url = %#v", d.Author.Url)
	}
	if d.Author.IconUrl == nil || *d.Author.IconUrl != "https://example.com/icon.png" {
		t.Fatalf("Author.IconUrl = %#v", d.Author.IconUrl)
	}
}

func TestFormConfig_DiscordNilAuthorAndTemplate(t *testing.T) {
	t.Parallel()

	src := &fc.FormConfig{
		FormId: 1,
		Notifiers: fc.NotifiersConfig{
			Discord: &discord.DiscordConfig{
				Webhooks: []*discord.WebhookKey{},
			},
		},
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	d := got.Notifiers.Discord
	if d == nil {
		t.Fatal("expected Discord config")
	}
	if d.Author != nil {
		t.Fatalf("Author = %#v, want nil", d.Author)
	}
	if d.Template != nil {
		t.Fatalf("Template = %#v, want nil", d.Template)
	}
}

func TestFormConfig_DiscordAuthorEmptyIcon(t *testing.T) {
	t.Parallel()

	src := &fc.FormConfig{
		FormId: 1,
		Notifiers: fc.NotifiersConfig{
			Discord: &discord.DiscordConfig{
				Author: &discord.Author{Name: "NoIcon"},
			},
		},
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	author := got.Notifiers.Discord.Author
	if author == nil {
		t.Fatal("expected Author")
	}
	if author.Name == nil || *author.Name != "NoIcon" {
		t.Fatalf("Name = %#v", author.Name)
	}
	if author.IconUrl != nil {
		t.Fatalf("IconUrl = %#v, want nil", author.IconUrl)
	}
}

func TestFormConfig_Full(t *testing.T) {
	t.Parallel()

	brevoTmpl, err := common.NewTemplate("email body {{email}}")
	if err != nil {
		t.Fatal(err)
	}
	discordTmpl, err := common.NewTemplate("discord {{email}}")
	if err != nil {
		t.Fatal(err)
	}

	src := &fc.FormConfig{
		FormId:        99,
		SiteId:        11,
		CaptchaType:   common.CaptchaTypeRecaptcha,
		CaptchaField:  "g-recaptcha-response",
		RedirectTo:    "/thanks",
		SchemaVersion: 2,
		FieldSchema: fc.FieldSchema{
			Fields: []fc.Field{
				{Name: "email", Type: fc.FieldTypeString, Required: true},
				{
					Name:     "scores",
					Type:     fc.FieldTypeArray,
					Required: true,
					Items:    &fc.FieldItems{Type: fc.FieldTypeNumber},
				},
			},
		},
		Notifiers: fc.NotifiersConfig{
			Brevo: &brevo.BrevoConfig{
				Subject:  "Subject",
				Template: brevoTmpl,
				Recipients: []*brevo.EmailContact{
					{Name: "Ops", Address: "ops@example.com"},
				},
			},
			Discord: &discord.DiscordConfig{
				Title:    "Title",
				Color:    1,
				Template: discordTmpl,
				Webhooks: []*discord.WebhookKey{
					{Snowflake: "sf", Token: "tok"},
				},
				Author: &discord.Author{
					Name: "Author",
					Icon: "https://cdn.example/icon.png",
				},
			},
		},
	}

	got, err := mappers.FormConfig(src)
	if err != nil {
		t.Fatal(err)
	}

	if got.Id != 99 || got.SiteId != 11 || got.SchemaVersion != 2 {
		t.Fatalf("ids = id=%d site=%d schema=%d", got.Id, got.SiteId, got.SchemaVersion)
	}
	if got.CaptchaType != api.Recaptcha {
		t.Fatalf("CaptchaType = %q", got.CaptchaType)
	}
	if got.CaptchaField == nil || *got.CaptchaField != "g-recaptcha-response" {
		t.Fatalf("CaptchaField = %#v", got.CaptchaField)
	}
	if got.RedirectTo == nil || *got.RedirectTo != "/thanks" {
		t.Fatalf("RedirectTo = %#v", got.RedirectTo)
	}
	if len(got.FieldSchema.Fields) != 2 {
		t.Fatalf("fields = %#v", got.FieldSchema.Fields)
	}
	if got.Notifiers.Brevo == nil || got.Notifiers.Discord == nil {
		t.Fatalf("notifiers = %#v", got.Notifiers)
	}
	if got.Notifiers.Brevo.Subject != "Subject" {
		t.Fatalf("brevo subject = %q", got.Notifiers.Brevo.Subject)
	}
	if got.Notifiers.Discord.Title == nil || *got.Notifiers.Discord.Title != "Title" {
		t.Fatalf("discord title = %#v", got.Notifiers.Discord.Title)
	}
}
