package httpserver

import (
	"context"

	m "github.com/imcrazytwkr/formdrain/models/http"
)

type contentTypeCtxKey struct{}

// WithContentType stores the parsed request body content type on ctx.
func WithContentType(ctx context.Context, contentType m.ContentType) context.Context {
	return context.WithValue(ctx, contentTypeCtxKey{}, contentType)
}
