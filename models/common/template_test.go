package common_test

import (
	"encoding/json"
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
)

func TestNewTemplate_Execute(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("Hello {{name}}")
	if err != nil {
		t.Fatal(err)
	}
	got := tmpl.ExecuteString(map[string]any{"name": "Ada"})
	if got != "Hello Ada" {
		t.Fatalf("got %q", got)
	}
}

func TestTemplate_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("x={{v}}")
	if err != nil {
		t.Fatal(err)
	}

	raw, err := json.Marshal(tmpl)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != `"x={{v}}"` {
		t.Fatalf("marshal = %s", raw)
	}

	var got common.Template
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if got.ExecuteString(map[string]any{"v": "1"}) != "x=1" {
		t.Fatalf("execute after unmarshal failed")
	}
}

func TestTemplate_UnmarshalInvalidJSON(t *testing.T) {
	t.Parallel()

	var got common.Template
	if err := json.Unmarshal([]byte(`42`), &got); err == nil {
		t.Fatal("expected error")
	}
}
