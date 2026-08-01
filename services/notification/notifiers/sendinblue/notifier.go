package sendinblue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/imcrazytwkr/formdrain/models/form_config/sendinblue"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers"
	"github.com/imcrazytwkr/formdrain/services/notification/notifiers/sendinblue/models"
	"github.com/imcrazytwkr/formdrain/utils/httpclient"
	"github.com/imcrazytwkr/formdrain/utils/httpclient/transports"
)

type sendinblueNotifier struct {
	sender *sendinblue.EmailContact
	client *http.Client
}

func NewSendinblueNotifier(
	senderName string,
	senderEmail string,
	client *http.Client,
) notifiers.SendinblueNotifier {
	return &sendinblueNotifier{
		sender: &sendinblue.EmailContact{
			Name:    senderName,
			Address: senderEmail,
		},
		client: httpclient.WithTransport(client, transports.LimitedTransport(client.Transport, rateLimiter)),
	}
}

func (n *sendinblueNotifier) send(request *models.Request) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, backendUrl, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	response, err := n.client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("sendinblue backend responded with %d", response.StatusCode)
	}

	return nil
}

func (n *sendinblueNotifier) Send(config *sendinblue.SendinblueConfig, form map[string]any) error {
	if len(config.Recipients) < 1 {
		return nil
	}

	return n.send(&models.Request{
		Sender:     n.sender,
		Recipients: config.Recipients,
		Subject:    config.Subject,
		Content:    config.RenderContent(form),
	})
}
