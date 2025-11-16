package utf8util

import (
	"bytes"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

// @NOTE: this whole package exists solely because some services
// *cought* recaptcha *cough* tend to set UTF-8 BOM mark in their
// responses at random
func FixBytes(body []byte) ([]byte, error) {
	return detectEncoding(body).NewDecoder().Bytes(body)
}

var utf8prefix = []byte{0xEF, 0xBB, 0xBF}
var utf16beprefix = []byte{0xFE, 0xFF}
var utf16leprefix = []byte{0xFF, 0xFE}

func detectEncoding(body []byte) encoding.Encoding {
	if len(body) < 2 {
		return encoding.Nop
	}

	if bytes.HasPrefix(body, utf8prefix) {
		return unicode.UTF8BOM
	}

	// @NOTE: using ExpectBOM because BOM is guaranteed to exist
	if bytes.HasPrefix(body, utf16beprefix) {
		return unicode.UTF16(unicode.BigEndian, unicode.ExpectBOM)
	}

	if bytes.HasPrefix(body, utf16leprefix) {
		return unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM)
	}

	return encoding.Nop
}
