package bodyparser

import (
	"errors"
	"io"
	"net/http"

	"github.com/imcrazytwkr/formdrain/httpserver/contenttype"
)

func Parse(req *http.Request) (map[string]any, error) {
	contentType := contenttype.GetContentType(req)
	parser, ok := Parsers[contentType]
	if !ok {
		return nil, ErrUnsupportedContentType
	}

	if req.ContentLength > maxBodySize {
		return nil, ErrBodyTooLarge
	}

	body, err := io.ReadAll(io.LimitReader(req.Body, maxReadSize))
	if err != nil {
		return nil, ErrMalformedBody
	}

	if len(body) > maxBodySize {
		return nil, ErrBodyTooLarge
	}

	data, err := parser.Parse(body)
	if err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			err = ErrMalformedBody
		}

		return nil, err
	}

	return data, nil
}
