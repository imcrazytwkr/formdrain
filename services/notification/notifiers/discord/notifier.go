package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/imcrazytwkr/formdrain/models/form_config/discord"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers/discord/models"
	"github.com/imcrazytwkr/formdrain/utils/httpclient"
	"github.com/imcrazytwkr/formdrain/utils/httpclient/transports"
)

type discordNotifier struct {
	userName string
	avatar   string
	client   *http.Client
}

func NewDiscordNotifier(userName string, avatar string, client *http.Client) notifiers.DiscordNotifier {
	return &discordNotifier{
		userName: userName,
		avatar:   avatar,
		client:   httpclient.WithTransport(client, transports.LimitedTransport(client.Transport, rateLimiter)),
	}
}

func (n *discordNotifier) makeRequest(embed *models.Embed) *models.Request {
	return &models.Request{
		Username: n.userName,
		Avatar:   n.avatar,
		Embeds:   []*models.Embed{embed},
	}
}

func (n *discordNotifier) send(snowflake string, token string, request *models.Request) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", snowflake, token),
		bytes.NewReader(payload),
	)
	if err != nil {
		return err
	}

	response, err := n.client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("discord backend responded with %d", response.StatusCode)
	}

	return nil
}

func (n *discordNotifier) Send(config *discord.DiscordConfig, form map[string]any) error {
	if len(config.Webhooks) < 1 {
		return nil
	}

	description, err := config.RenderContent(form)
	if err != nil {
		return err
	}

	request := n.makeRequest(&models.Embed{
		Author:      config.Author,
		Title:       config.Title,
		Url:         config.Url,
		Description: description,
		Color:       config.Color,
	})

	var errs []error
	for _, key := range config.Webhooks {
		err := n.send(key.Snowflake, key.Token, request)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
