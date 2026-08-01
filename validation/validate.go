package validation

import (
	"errors"
	"fmt"
	"strconv"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

// FieldError ties a validation failure to a form field name.
type FieldError struct {
	Field string
	Err   error
}

func (e *FieldError) Error() string {
	if e == nil || e.Err == nil {
		return "form field validation failed"
	}
	if len(e.Field) < 1 {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Err.Error())
}

func (e *FieldError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func fieldErr(name string, err error) error {
	return &FieldError{Field: name, Err: err}
}

// ValidateFormPayload checks data against schema and returns a cleaned payload
// suitable for persistence (coerced types, schema keys only).
// On failure it returns errors.Join of all field errors (no partial payload).
func ValidateFormPayload(schema fc.FieldSchema, data map[string]any) (map[string]any, error) {
	allowed := make(map[string]fc.Field, len(schema.Fields))
	for _, field := range schema.Fields {
		allowed[field.Name] = field
	}

	var errs []error

	for key := range data {
		if _, ok := allowed[key]; !ok {
			errs = append(errs, fieldErr(key, ErrUnknownField))
		}
	}

	out := make(map[string]any, len(schema.Fields))
	for _, field := range schema.Fields {
		raw, present := data[field.Name]
		if !present || raw == nil {
			if field.Required {
				errs = append(errs, fieldErr(field.Name, ErrMissingRequiredField))
			}
			continue
		}

		value, err := normalizeValue(field, raw)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if value == nil {
			if field.Required {
				errs = append(errs, fieldErr(field.Name, ErrMissingRequiredField))
			}
			continue
		}

		out[field.Name] = value
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return out, nil
}

func normalizeValue(field fc.Field, raw any) (any, error) {
	switch field.Type {
	case fc.FieldTypeString:
		return normalizeString(field.Name, raw)
	case fc.FieldTypeNumber:
		return normalizeNumber(field.Name, raw)
	case fc.FieldTypeBoolean:
		return normalizeBoolean(field.Name, raw)
	case fc.FieldTypeArray:
		return normalizeArray(field, raw)
	default:
		return nil, fieldErr(field.Name, ErrUnsupportedFieldType)
	}
}

func normalizeString(name string, raw any) (any, error) {
	switch v := raw.(type) {
	case string:
		return v, nil
	default:
		return nil, fieldErr(name, ErrInvalidFieldType)
	}
}

func normalizeNumber(name string, raw any) (any, error) {
	switch v := raw.(type) {
	case int64:
		return v, nil
	case float64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i, nil
		}
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fieldErr(name, ErrInvalidFieldType)
		}
		return f, nil
	default:
		return nil, fieldErr(name, ErrInvalidFieldType)
	}
}

func normalizeBoolean(name string, raw any) (any, error) {
	switch v := raw.(type) {
	case bool:
		return v, nil
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fieldErr(name, ErrInvalidFieldType)
		}
		return b, nil
	default:
		return nil, fieldErr(name, ErrInvalidFieldType)
	}
}

func normalizeArray(field fc.Field, raw any) (any, error) {
	if field.Items == nil {
		return nil, fieldErr(field.Name, ErrMissingArrayItemType)
	}

	itemType := field.Items.Type
	switch itemType {
	case fc.FieldTypeString, fc.FieldTypeNumber, fc.FieldTypeBoolean:
		// ok
	default:
		return nil, fieldErr(field.Name, ErrUnsupportedItemType)
	}

	elements, err := asArrayElements(field.Name, raw)
	if err != nil {
		return nil, err
	}

	out := make([]any, 0, len(elements))
	for _, el := range elements {
		if el == nil {
			return nil, fieldErr(field.Name, ErrInvalidArrayItems)
		}

		var normalized any
		switch itemType {
		case fc.FieldTypeString:
			normalized, err = normalizeString(field.Name, el)
		case fc.FieldTypeNumber:
			normalized, err = normalizeNumber(field.Name, el)
		case fc.FieldTypeBoolean:
			normalized, err = normalizeBoolean(field.Name, el)
		}
		if err != nil {
			return nil, fieldErr(field.Name, ErrInvalidArrayItems)
		}
		out = append(out, normalized)
	}

	return out, nil
}

func asArrayElements(name string, raw any) ([]any, error) {
	switch v := raw.(type) {
	case []any:
		return v, nil
	case []string:
		out := make([]any, len(v))
		for i, s := range v {
			out[i] = s
		}
		return out, nil
	case string, bool, int, int64, float64:
		return []any{v}, nil
	default:
		return nil, fieldErr(name, ErrInvalidFieldType)
	}
}
