package safemustache

import (
	"io"

	"github.com/cbroglie/mustache"
	"github.com/imcrazytwkr/formdrain/templates"
	"github.com/imcrazytwkr/formdrain/templates/common"
)

type safeTemplate struct {
	inner *mustache.Template
}

func NewTemplate(inner *mustache.Template) (templates.Template, error) {
	if inner == nil {
		return nil, common.ErrNilTemplate
	}

	return &safeTemplate{inner}, nil
}

func (t *safeTemplate) Execute(w io.Writer, data map[string]any) error {
	if t == nil || t.inner == nil {
		return common.ErrNilTemplate
	}

	return t.inner.FRender(w, data)
}

func (t *safeTemplate) ExecuteString(data map[string]any) (string, error) {
	if t == nil || t.inner == nil {
		return "", common.ErrNilTemplate
	}

	return t.inner.Render(data)
}
