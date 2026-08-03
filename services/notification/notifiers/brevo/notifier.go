package brevo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/imcrazytwkr/formdrain/constants"
	"github.com/imcrazytwkr/formdrain/models/form_config/brevo"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers/brevo/models"
	"github.com/imcrazytwkr/formdrain/utils/httpclient"
	"github.com/imcrazytwkr/formdrain/utils/httpclient/transports"
)

type brevoNotifier struct {
	sender *brevo.EmailContact
	apiKey string
	client *http.Client
}

func NewBrevoNotifier(
	senderName string,
	senderEmail string,
	apiKey string,
	client *http.Client,
) notifiers.BrevoNotifier {
	return &brevoNotifier{
		sender: &brevo.EmailContact{
			Name:    senderName,
			Address: senderEmail,
		},
		apiKey: apiKey,
		client: httpclient.WithTransport(client, transports.LimitedTransport(client.Transport, rateLimiter)),
	}
}

func (n *brevoNotifier) send(ctx context.Context, request *models.Request) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, backendUrl, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set(constants.HeaderAccept, constants.ContentTypeJson)
	req.Header.Set(constants.HeaderContentType, constants.ContentTypeJson)
	req.Header.Set(headerApiKey, n.apiKey)

	response, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusCreated, http.StatusOK:
		return nil
	default:
		return fmt.Errorf("brevo backend responded with %d", response.StatusCode)
	}
}

func (n *brevoNotifier) Send(ctx context.Context, config *brevo.BrevoConfig, form map[string]any) error {
	if len(config.Recipients) < 1 {
		return nil
	}

	content, err := config.RenderContent(form)
	if err != nil {
		return err
	}

	return n.send(ctx, &models.Request{
		Sender:     n.sender,
		Recipients: config.Recipients,
		Subject:    config.Subject,
		Content:    content,
	})
}
