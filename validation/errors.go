package validation

import "errors"

const errValidationFailed = "form validation failed"

var (
	ErrUnknownField         = errors.New("unknown form field")
	ErrMissingRequiredField = errors.New("required form field is missing")
	ErrInvalidFieldType     = errors.New("form field has invalid type")
	ErrInvalidArrayItems    = errors.New("form field array items are invalid")
	ErrMissingArrayItemType = errors.New("array field schema is missing items.type")
	ErrUnsupportedFieldType = errors.New("unsupported field type in schema")
	ErrUnsupportedItemType  = errors.New("unsupported array item type in schema")
	ErrValidationFailed     = errors.New(errValidationFailed)
)
