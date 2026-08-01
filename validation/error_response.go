package validation

import (
	"errors"

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
