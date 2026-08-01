package maputil_test

import (
	"testing"

	"github.com/imcrazytwkr/formdrain/utils/maputil"
)

func TestGetString(t *testing.T) {
	t.Parallel()

	obj := map[string]any{
		"name":  "ada",
		"count": 1,
		"empty": "",
	}

	if got, ok := maputil.GetString(obj, "name"); !ok || got != "ada" {
		t.Fatalf("name: got (%q, %v)", got, ok)
	}
	if got, ok := maputil.GetString(obj, "missing"); ok || got != "" {
		t.Fatalf("missing: got (%q, %v)", got, ok)
	}
	if got, ok := maputil.GetString(obj, "count"); ok || got != "" {
		t.Fatalf("wrong type: got (%q, %v)", got, ok)
	}
	if got, ok := maputil.GetString(obj, "empty"); !ok || got != "" {
		t.Fatalf("empty: got (%q, %v)", got, ok)
	}
}
