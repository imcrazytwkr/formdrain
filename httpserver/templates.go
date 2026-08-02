package httpserver

import (
	"fmt"

	"github.com/imcrazytwkr/formdrain/templates"
)

// GetTemplate returns the embedded template for name.
// Panics if name is not a known embedded template (programmer error).
func GetTemplate(name string) templates.Template {
	tmpl, ok := embeddedTemplates[name]
	if !ok || tmpl == nil {
		panic(fmt.Sprintf("httpserver: unknown embedded template %q", name))
	}

	return tmpl
}
