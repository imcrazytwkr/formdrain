package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/imcrazytwkr/formdrain/httpserver"
)

func TestEmbeddedTemplatesRender(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		template string
		data     any
		want     string
	}{
		{
			name:     "success",
			template: "form/success",
			data:     nil,
			want:     "Form submitted",
		},
		{
			name:     "redirect",
			template: "form/redirect",
			data:     nil,
			want:     "Redirecting",
		},
		{
			name:     "generic error",
			template: "errors/generic",
			data: map[string]any{
				"status":  404,
				"title":   "Not Found",
				"message": "gone",
			},
			want: "Not Found",
		},
		{
			name:     "validation error",
			template: "errors/validation",
			data: map[string]any{
				"status":  400,
				"message": "form validation failed",
				"errors": []map[string]string{
					{"field": "email", "message": "required"},
				},
			},
			want: "email",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			httpserver.HandleResponse(t.Context(), w, http.StatusOK, tc.template, tc.data)
			if !strings.Contains(w.Body.String(), tc.want) {
				t.Fatalf("body = %q, want substring %q", w.Body.String(), tc.want)
			}
		})
	}
}
