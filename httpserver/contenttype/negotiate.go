package contenttype

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
	"github.com/imcrazytwkr/formdrain/utils/stringutil"
)

// Negotiate picks the best Content-Type from offers based on the Accept header.
// Returns ContentTypeUndefined when Accept is empty or nothing matches — no default.
func Negotiate(r *http.Request, offers []m.ContentType) m.ContentType {
	accept := strings.TrimSpace(r.Header.Get(constants.HeaderAccept))
	if len(accept) < 1 {
		return m.ContentTypeUndefined
	}

	type candidate struct {
		offer m.ContentType
		q     float64
	}

	var best candidate
	best.q = -1

	for part := range strings.SplitSeq(accept, ",") {
		part = strings.TrimSpace(part)
		if len(part) < 1 {
			continue
		}

		media := part
		q := 1.0
		if before, params, ok := strings.Cut(part, ";"); ok {
			media = strings.TrimSpace(before)
			for p := range strings.SplitSeq(params, ";") {
				v, ok := strings.CutPrefix(strings.TrimSpace(p), "q=")
				if !ok {
					continue
				}

				parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
				if err == nil {
					q = parsed
				}
			}
		}

		if q <= 0 {
			continue
		}

		for _, offer := range offers {
			if offer == m.ContentTypeUndefined {
				continue
			}

			if !acceptMatches(media, offer.String()) {
				continue
			}

			if q > best.q {
				best.offer = offer
				best.q = q
			}

			break
		}
	}

	if best.q < 0 {
		return m.ContentTypeUndefined
	}

	return best.offer
}

func acceptMatches(acceptMedia, offer string) bool {
	if acceptMedia == "*/*" || acceptMedia == offer {
		return true
	}

	// type/*
	prefix, ok := strings.CutSuffix(acceptMedia, "/*")
	if !ok {
		return false
	}

	offerType := stringutil.TakeUntilByte(offer, '/')
	return prefix == offerType
}
