package common

import (
	"github.com/valyala/fasttemplate"
	"go.mongodb.org/mongo-driver/bson"
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

func (t *Template) MarshalBSON() ([]byte, error) {
	return bson.Marshal(t.source)
}

func (t *Template) UnmarshalBSON(b []byte) error {
	var raw string
	err := bson.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	t.source = raw
	return t.Reset(raw, TemplateStartTag, TemplateEndTag)
}
