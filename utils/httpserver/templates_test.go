package httpserver_test

import (
	"path/filepath"
	"testing"

	"github.com/imcrazytwkr/formdrain/utils/httpserver"
)

func TestLoadTemplates(t *testing.T) {
	dir := filepath.Join("..", "..", "templates")
	err := httpserver.LoadTemplatesFromPath(dir)
	if err != nil {
		t.Fatal(err)
	}
}
