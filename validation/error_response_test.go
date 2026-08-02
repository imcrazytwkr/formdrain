package validation

import (
	"net/http"
	"testing"
)

func TestHTMLContext(t *testing.T) {
	t.Parallel()

	t.Run("empty errors", func(t *testing.T) {
		t.Parallel()
		res := &ValidationErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "form validation failed",
		}

		got := res.TemplateData()
		if got["status"] != http.StatusBadRequest || got["message"] != "form validation failed" {
			t.Fatalf("got = %#v", got)
		}
		errors, ok := got["errors"].([]map[string]string)
		if !ok || len(errors) != 0 {
			t.Fatalf("errors = %#v", got["errors"])
		}
	})

	t.Run("sorted field entries", func(t *testing.T) {
		t.Parallel()
		res := &ValidationErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "form validation failed",
			Errors: map[string]string{
				"zeta":  "z",
				"alpha": "a",
			},
		}

		got := res.TemplateData()
		errors, ok := got["errors"].([]map[string]string)
		if !ok || len(errors) != 2 {
			t.Fatalf("errors = %#v", got["errors"])
		}
		if errors[0]["field"] != "alpha" || errors[0]["message"] != "a" {
			t.Fatalf("first = %#v", errors[0])
		}
		if errors[1]["field"] != "zeta" || errors[1]["message"] != "z" {
			t.Fatalf("second = %#v", errors[1])
		}
	})
}
