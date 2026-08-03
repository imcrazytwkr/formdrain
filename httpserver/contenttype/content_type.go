package contenttype

import (
	"mime"
	"net/http"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
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
