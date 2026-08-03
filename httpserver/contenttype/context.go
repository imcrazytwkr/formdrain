package contenttype

import (
	"context"
	"net/http"

	m "github.com/imcrazytwkr/formdrain/models/http"
)

type contentTypeCtxKey struct{}

// WithContentType stores the parsed request body content type on ctx.
func WithContentType(ctx context.Context, contentType m.ContentType) context.Context {
	return context.WithValue(ctx, contentTypeCtxKey{}, contentType)
}

// GetContentType returns the parsed body content type stored on the request context.
func GetContentType(r *http.Request) m.ContentType {
	raw, ok := r.Context().Value(contentTypeCtxKey{}).(m.ContentType)
	if !ok {
		return m.ContentTypeUndefined
	}

	return raw
}
