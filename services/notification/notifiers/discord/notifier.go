package discord

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
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

func (n *discordNotifier) send(ctx context.Context, snowflake string, token string, request *models.Request) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", snowflake, token),
		bytes.NewReader(payload),
	)
	if err != nil {
		return err
	}

	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJson)

	response, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusNoContent, http.StatusOK:
		return nil
	default:
		return fmt.Errorf("discord backend responded with %d", response.StatusCode)
	}
}

func (n *discordNotifier) Send(ctx context.Context, config *discord.DiscordConfig, form map[string]any) error {
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
		err := n.send(ctx, key.Snowflake, key.Token, request)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
