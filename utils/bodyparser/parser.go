package bodyparser

import (
	"github.com/imcrazytwkr/formdrain/models/http"
)

type BodyParser interface {
	Parse(body []byte) (map[string]any, error)
}

var Parsers = map[http.ContentType]BodyParser{
	http.ContentTypeHTML: NewFormParser(),
	http.ContentTypeJSON: NewJsonParser(),
}
