package safemustache_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/cbroglie/mustache"
	"github.com/imcrazytwkr/formdrain/utils/safemustache"
)

func TestAllowMissingVariablesDisabled(t *testing.T) {
	if mustache.AllowMissingVariables {
		t.Fatal("expected AllowMissingVariables to be false after importing safemustache")
	}
}

func TestParse_Rejects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		src     string
		wantErr error
	}{
		{name: "empty", src: "", wantErr: safemustache.ErrEmpty},
		{name: "raw triple", src: "Hello {{{name}}}", wantErr: safemustache.ErrRawInterpolation},
		{name: "raw amp", src: "Hello {{& name}}", wantErr: safemustache.ErrRawInterpolation},
		{name: "raw amp tight", src: "Hello {{&name}}", wantErr: safemustache.ErrRawInterpolation},
		{name: "delimiter change", src: "{{=<% %>=}}<%name%>", wantErr: safemustache.ErrDelimiterChange},
		{name: "partial", src: "x {{> foo}} y", wantErr: safemustache.ErrPartial},
		{
			name:    "too large",
			src:     strings.Repeat("a", safemustache.MaxSourceBytes+1),
			wantErr: safemustache.ErrTooLarge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := safemustache.Parse(tt.src)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestParse_HappyAndEscape(t *testing.T) {
	t.Parallel()

	tmpl, err := safemustache.Parse("Hello {{name}}")
	if err != nil {
		t.Fatal(err)
	}

	got, err := tmpl.Render(map[string]any{"name": "world"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "Hello world" {
		t.Fatalf("got %q", got)
	}

	got, err = tmpl.Render(map[string]any{"name": "<script>x</script>"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(got, "<script>") {
		t.Fatalf("expected escaped output, got %q", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Fatalf("expected html entities, got %q", got)
	}
}

func TestParse_SectionsAndLists(t *testing.T) {
	t.Parallel()

	tmpl, err := safemustache.Parse("{{#items}}[{{.}}]{{/items}}")
	if err != nil {
		t.Fatal(err)
	}

	got, err := tmpl.Render(map[string]any{"items": []string{"a", "b"}})
	if err != nil {
		t.Fatal(err)
	}
	if got != "[a][b]" {
		t.Fatalf("got %q", got)
	}

	got, err = tmpl.Render(map[string]any{"items": []string{}})
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("empty list: got %q", got)
	}
}

func TestParse_ErrorsEntryList(t *testing.T) {
	t.Parallel()

	src := `{{#errors}}{{field}}:{{message}};{{/errors}}{{^errors}}none{{/errors}}`
	tmpl, err := safemustache.Parse(src)
	if err != nil {
		t.Fatal(err)
	}

	got, err := tmpl.Render(map[string]any{
		"errors": []map[string]string{
			{"field": "email", "message": "required"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != "email:required;" {
		t.Fatalf("got %q", got)
	}

	got, err = tmpl.Render(map[string]any{"errors": []map[string]string{}})
	if err != nil {
		t.Fatal(err)
	}
	if got != "none" {
		t.Fatalf("inverted empty: got %q", got)
	}
}

func TestParse_DottedNames(t *testing.T) {
	t.Parallel()

	tmpl, err := safemustache.Parse("{{user.name}}")
	if err != nil {
		t.Fatal(err)
	}

	got, err := tmpl.Render(map[string]any{
		"user": map[string]any{"name": "Ada"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != "Ada" {
		t.Fatalf("got %q", got)
	}
}

func TestParse_MissingVariable(t *testing.T) {
	t.Parallel()

	tmpl, err := safemustache.Parse("Hello {{name}}")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmpl.Render(map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestParse_MaxExactSizeOK(t *testing.T) {
	t.Parallel()

	prefix := strings.Repeat("a", safemustache.MaxSourceBytes-5)
	src := prefix + "{{x}}"
	if len(src) != safemustache.MaxSourceBytes {
		t.Fatalf("len = %d", len(src))
	}
	tmpl, err := safemustache.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	got, err := tmpl.Render(map[string]any{"x": "z"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(got, "z") {
		t.Fatalf("got %q", got)
	}
}
