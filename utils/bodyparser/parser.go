package bodyparser

import (
	"github.com/imcrazytwkr/formdrain/constants"
	m "github.com/imcrazytwkr/formdrain/models/http"
)

type BodyParser interface {
	Parse(body []byte) (map[string]any, error)
}

var Parsers = map[string]BodyParser{
	constants.ContentTypeForm: NewFormParser(),
	constants.ContentTypeJson: NewJsonParser(),
}

var ParsersNew = map[m.ContentType]BodyParser{
	m.ContentTypeHTML: NewFormParser(),
	m.ContentTypeJSON: NewJsonParser(),
}
