package httpclient

import (
	"net/http"

	"github.com/imcrazytwkr/formdrain/types"
	"go.uber.org/ratelimit"
)

type limitedHttpClient struct {
	client types.HttpClient
	limit  ratelimit.Limiter
}

func NewLimitedHttpClient(client types.HttpClient, rate int, opts ...ratelimit.Option) types.HttpClient {
	return &limitedHttpClient{
		client: client,
		limit:  ratelimit.New(rate, opts...),
	}
}

func (c *limitedHttpClient) Do(request *http.Request) (*http.Response, error) {
	c.limit.Take()
	return c.client.Do(request)
}
