package common

import (
	"encoding/json"
	"errors"

	"github.com/cbroglie/mustache"
	"github.com/imcrazytwkr/formdrain/utils/safemustache"
)

type Template struct {
	source string
	inner  *mustache.Template
}

// @api: internal
var ErrNilTemplate = errors.New("common: nil template")

func NewTemplate(raw string) (*Template, error) {
	inner, err := safemustache.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &Template{
		source: raw,
		inner:  inner,
	}, nil
}

func (t *Template) ExecuteString(data map[string]any) (string, error) {
	if t == nil || t.inner == nil {
		return "", ErrNilTemplate
	}

	return t.inner.Render(data)
}

func (t Template) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.source)
}

func (t *Template) UnmarshalJSON(b []byte) error {
	var raw string
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	tmpl, err := NewTemplate(raw)
	if err != nil {
		return err
	}

	*t = *tmpl
	return nil
}
