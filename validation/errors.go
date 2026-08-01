package validation

import "errors"

const errValidationFailed = "form validation failed"

var ErrUnknownField = errors.New("unknown form field")
var ErrMissingRequiredField = errors.New("required form field is missing")
var ErrInvalidFieldType = errors.New("form field has invalid type")
var ErrInvalidArrayItems = errors.New("form field array items are invalid")
var ErrMissingArrayItemType = errors.New("array field schema is missing items.type")
var ErrUnsupportedFieldType = errors.New("unsupported field type in schema")
var ErrUnsupportedItemType = errors.New("unsupported array item type in schema")
var ErrValidationFailed = errors.New(errValidationFailed)
