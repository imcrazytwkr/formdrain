package common

import (
	"encoding/json"

	"github.com/valyala/fasttemplate"
)

type Template struct {
	source string
	*fasttemplate.Template
}

const TemplateStartTag = "{{"
const TemplateEndTag = "}}"

func NewTemplate(raw string) (*Template, error) {
	tmpl, err := fasttemplate.NewTemplate(raw, TemplateStartTag, TemplateEndTag)
	if err != nil {
		return nil, err
	}

	return &Template{
		source:   raw,
		Template: tmpl,
	}, nil
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
