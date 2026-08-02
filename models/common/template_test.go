package common_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/imcrazytwkr/formdrain/models/common"
	"github.com/imcrazytwkr/formdrain/utils/safemustache"
)

func TestNewTemplate_Execute(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("Hello {{name}}")
	if err != nil {
		t.Fatal(err)
	}
	got, err := tmpl.ExecuteString(map[string]any{"name": "Ada"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "Hello Ada" {
		t.Fatalf("got %q", got)
	}
}

func TestTemplate_EscapesHTML(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("{{v}}")
	if err != nil {
		t.Fatal(err)
	}
	got, err := tmpl.ExecuteString(map[string]any{"v": "<x>"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "&lt;x&gt;" {
		t.Fatalf("got %q", got)
	}
}

func TestTemplate_MissingVariable(t *testing.T) {
	t.Parallel()

	tmpl, err := common.NewTemplate("Hello {{name}}")
	if err != nil {
		t.Fatal(err)
	}
	_, err = tmpl.ExecuteString(map[string]any{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTemplate_RejectsRaw(t *testing.T) {
	t.Parallel()

	_, err := common.NewTemplate("{{{name}}}")
	if !errors.Is(err, safemustache.ErrRawInterpolation) {
		t.Fatalf("err = %v", err)
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
	out, err := got.ExecuteString(map[string]any{"v": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "x=1" {
		t.Fatalf("got %q", out)
	}
}

func TestTemplate_UnmarshalInvalidJSON(t *testing.T) {
	t.Parallel()

	var got common.Template
	if err := json.Unmarshal([]byte(`42`), &got); err == nil {
		t.Fatal("expected error")
	}
}

func TestTemplate_NilExecute(t *testing.T) {
	t.Parallel()

	var tmpl *common.Template
	_, err := tmpl.ExecuteString(nil)
	if !errors.Is(err, common.ErrNilTemplate) {
		t.Fatalf("err = %v", err)
	}
}
