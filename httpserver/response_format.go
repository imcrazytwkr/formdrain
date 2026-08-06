package httpserver

import (
	"mime"
	"net/http"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/stringutil"
)

// ResponseFormat reads the response Content-Type header.
// If unset or unrecognized, falls back to HTML.
func ResponseFormat(w http.ResponseWriter) m.ContentType {
	raw := w.Header().Get(constants.HeaderContentType)
	if len(raw) < 1 {
		return m.ContentTypeHTML
	}

	mediaType, _, err := mime.ParseMediaType(raw)
	if err != nil {
		mediaType = strings.TrimSpace(stringutil.TakeUntilByte(raw, ';'))
	}

	format := m.ParseContentType(mediaType)
	if format == m.ContentTypeUndefined {
		return m.ContentTypeHTML
	}

	return format
}
