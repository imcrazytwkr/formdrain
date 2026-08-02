package httpserver

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/imcrazytwkr/formdrain/templates"
	"github.com/imcrazytwkr/formdrain/templates/safemustache"
)

const templatesDir = "templates"

//go:embed templates
var templateFS embed.FS

var embeddedTemplates = make(map[string]templates.Template)

func init() {
	err := fs.WalkDir(templateFS, templatesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".html") {
			return nil
		}

		data, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(strings.TrimPrefix(filepath.ToSlash(path), "templates/"), ".html")
		tmpl, err := safemustache.Parse(string(data))
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}

		embeddedTemplates[name] = tmpl
		return nil
	})

	if err != nil {
		panic(fmt.Sprintf("httpserver: load embedded templates: %v", err))
	}

	if len(embeddedTemplates) < 1 {
		panic("httpserver: load embedded templates: no html templates embedded")
	}
}
