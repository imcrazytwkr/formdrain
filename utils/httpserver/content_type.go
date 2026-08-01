package httpserver

import (
	"mime"
	"net/http"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/stringutil"
)

// RequestContentType returns the request Content-Type MIME type without parameters.
func RequestContentType(r *http.Request) string {
	header := r.Header.Get(constants.HeaderContentType)

	mediaType, _, err := mime.ParseMediaType(header)
	if err == nil {
		return mediaType
	}

	return strings.TrimSpace(stringutil.TakeUntilByte(header, ';'))
}

// GetContentType returns the parsed body content type stored on the request context.
func GetContentType(r *http.Request) m.ContentType {
	raw, ok := r.Context().Value(contentTypeCtxKey{}).(m.ContentType)
	if !ok {
		return m.ContentTypeUndefined
	}

	return raw
}
