package httpserver

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const DefaultTemplateDirectory = "templates"

var templates *template.Template

func LoadTemplates() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	return LoadTemplatesFromPath(filepath.Join(filepath.Dir(exePath), DefaultTemplateDirectory))
}

// LoadTemplatesFromPath parses *.html files under dir. Template names are paths relative
// to dir using forward slashes (e.g. "errors/generic.html").
func LoadTemplatesFromPath(dir string) error {
	root := template.New("httpserver")
	var found bool

	// I am extremely disappointed this has to be done manually
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".html") {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		name := filepath.ToSlash(rel)
		_, err = root.New(name).Parse(string(data))
		if err != nil {
			return fmt.Errorf("parse template %q: %w", name, err)
		}

		found = true
		return nil
	})
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("no html templates found under %q", dir)
	}

	templates = root
	return nil
}
