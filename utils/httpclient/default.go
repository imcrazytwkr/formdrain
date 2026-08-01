package httpclient

import (
	"net/http"
	"time"

	"github.com/imcrazytwkr/formdrain/utils/httpclient/transports"
)

const defaultRequestTimeout = 120 * time.Second

func DefaultClient() *http.Client {
	return &http.Client{
		Transport: transports.DefaultTransport(),
		Timeout:   defaultRequestTimeout,
	}
}
