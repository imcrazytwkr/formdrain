package common

import (
	"encoding/json"
	"io"

	"github.com/imcrazytwkr/formdrain/templates"
	"github.com/imcrazytwkr/formdrain/templates/common"
	"github.com/imcrazytwkr/formdrain/templates/safemustache"
)

type Template struct {
	source string
	inner  templates.Template
}

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

func (t *Template) Execute(w io.Writer, data map[string]any) error {
	if t == nil || t.inner == nil {
		return common.ErrNilTemplate
	}

	return t.inner.Execute(w, data)
}

func (t *Template) ExecuteString(data map[string]any) (string, error) {
	if t == nil || t.inner == nil {
		return "", common.ErrNilTemplate
	}

	return t.inner.ExecuteString(data)
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
