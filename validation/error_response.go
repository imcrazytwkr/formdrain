package validation

import (
	"errors"
	"slices"

	"github.com/imcrazytwkr/formdrain/utils/errorutil"
)

type ValidationErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// NewValidationErrorResponse builds the JSON error object for validation failures.
func NewValidationErrorResponse(status int, err error) *ValidationErrorResponse {
	response := &ValidationErrorResponse{
		Status: status,
	}

	errs, isMulti := errorutil.UnwrapMultiErr(err)
	if !isMulti || len(errs) == 0 {
		response.Message = err.Error()
		return response
	}

	response.Message = errValidationFailed
	response.Errors = FieldErrorsMap(errs)
	return response
}

func FieldErrorsMap(errs []error) map[string]string {
	// Just a sanity check
	if len(errs) == 0 {
		return nil
	}

	out := make(map[string]string)
	for _, err := range errs {
		var fieldErr *FieldError
		if !errors.As(err, &fieldErr) || fieldErr == nil || len(fieldErr.Field) == 0 {
			continue
		}

		if fieldErr.Err != nil {
			out[fieldErr.Field] = fieldErr.Err.Error()
			continue
		}

		out[fieldErr.Field] = fieldErr.Error()
	}

	return out
}

// TemplateData returns Mustache-friendly data for the validation error HTML template.
func (r *ValidationErrorResponse) TemplateData() map[string]any {
	if r == nil {
		return map[string]any{
			"status":  0,
			"message": "",
			"errors":  []map[string]string{},
		}
	}

	entries := make([]map[string]string, 0, len(r.Errors))
	fields := make([]string, 0, len(r.Errors))
	for field := range r.Errors {
		fields = append(fields, field)
	}

	slices.Sort(fields)
	for _, field := range fields {
		entries = append(entries, map[string]string{
			"field":   field,
			"message": r.Errors[field],
		})
	}

	return map[string]any{
		"status":  r.Status,
		"message": r.Message,
		"errors":  entries,
	}
}
