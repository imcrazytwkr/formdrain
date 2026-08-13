package validation_test

import (
	"errors"
	"net/http"
	"testing"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/validation"
)

func schema(fields ...fc.Field) fc.FieldSchema {
	return fc.FieldSchema{Version: 1, Fields: fields}
}

func TestValidateFormPayload_HappyPath(t *testing.T) {
	s := schema(
		fc.Field{Name: "email", Type: fc.FieldTypeString, Required: true},
		fc.Field{Name: "age", Type: fc.FieldTypeNumber, Required: false},
		fc.Field{Name: "subscribe", Type: fc.FieldTypeBoolean, Required: false},
		fc.Field{
			Name:     "interests",
			Type:     fc.FieldTypeArray,
			Required: false,
			Items:    &fc.FieldItems{Type: fc.FieldTypeString},
		},
	)

	got, err := validation.ValidateFormPayload(s, map[string]any{
		"email":     "a@b.c",
		"age":       int64(30),
		"subscribe": true,
		"interests": []any{"go", "sql"},
	})
	if err != nil {
		t.Fatalf("ValidateFormPayload: %v", err)
	}
	if got["email"] != "a@b.c" || got["age"] != int64(30) || got["subscribe"] != true {
		t.Fatalf("got %#v", got)
	}
	interests, ok := got["interests"].([]any)
	if !ok || len(interests) != 2 {
		t.Fatalf("interests: %#v", got["interests"])
	}
}

func TestValidateFormPayload_UnknownKey(t *testing.T) {
	s := schema(fc.Field{Name: "email", Type: fc.FieldTypeString, Required: true})
	_, err := validation.ValidateFormPayload(s, map[string]any{
		"email": "a@b.c",
		"extra": "nope",
	})
	if !errors.Is(err, validation.ErrUnknownField) {
		t.Fatalf("got %v", err)
	}
}

func TestValidateFormPayload_RequiredMissing(t *testing.T) {
	s := schema(fc.Field{Name: "email", Type: fc.FieldTypeString, Required: true})
	_, err := validation.ValidateFormPayload(s, map[string]any{})
	if !errors.Is(err, validation.ErrMissingRequiredField) {
		t.Fatalf("got %v", err)
	}
}

func TestValidateFormPayload_OptionalMissingOK(t *testing.T) {
	s := schema(fc.Field{Name: "nickname", Type: fc.FieldTypeString, Required: false})
	got, err := validation.ValidateFormPayload(s, map[string]any{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %#v", got)
	}
}

func TestValidateFormPayload_TypeCoercion(t *testing.T) {
	s := schema(
		fc.Field{Name: "age", Type: fc.FieldTypeNumber, Required: true},
		fc.Field{Name: "ok", Type: fc.FieldTypeBoolean, Required: true},
	)
	got, err := validation.ValidateFormPayload(s, map[string]any{
		"age": "42",
		"ok":  "true",
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got["age"] != int64(42) || got["ok"] != true {
		t.Fatalf("got %#v", got)
	}
}

func TestValidateFormPayload_InvalidType(t *testing.T) {
	s := schema(fc.Field{Name: "age", Type: fc.FieldTypeNumber, Required: true})
	_, err := validation.ValidateFormPayload(s, map[string]any{"age": "nope"})
	if !errors.Is(err, validation.ErrInvalidFieldType) {
		t.Fatalf("got %v", err)
	}
}

func TestValidateFormPayload_ArrayScalarCoerce(t *testing.T) {
	s := schema(fc.Field{
		Name:     "tags",
		Type:     fc.FieldTypeArray,
		Required: true,
		Items:    &fc.FieldItems{Type: fc.FieldTypeString},
	})
	got, err := validation.ValidateFormPayload(s, map[string]any{"tags": "solo"})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	tags, ok := got["tags"].([]any)
	if !ok || len(tags) != 1 || tags[0] != "solo" {
		t.Fatalf("tags: %#v", got["tags"])
	}
}

func TestValidateFormPayload_ArrayFromFormSlice(t *testing.T) {
	s := schema(fc.Field{
		Name:     "tags",
		Type:     fc.FieldTypeArray,
		Required: true,
		Items:    &fc.FieldItems{Type: fc.FieldTypeString},
	})
	got, err := validation.ValidateFormPayload(s, map[string]any{
		"tags": []string{"a", "b"},
	})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	tags := got["tags"].([]any)
	if len(tags) != 2 || tags[0] != "a" || tags[1] != "b" {
		t.Fatalf("tags: %#v", tags)
	}
}

func TestValidateFormPayload_RequiredEmptyArray(t *testing.T) {
	s := schema(fc.Field{
		Name:     "tags",
		Type:     fc.FieldTypeArray,
		Required: true,
		Items:    &fc.FieldItems{Type: fc.FieldTypeString},
	})

	_, err := validation.ValidateFormPayload(s, map[string]any{"tags": []any{}})
	if !errors.Is(err, validation.ErrMissingRequiredField) {
		t.Fatalf("[]any{}: %v", err)
	}

	_, err = validation.ValidateFormPayload(s, map[string]any{"tags": []string{}})
	if !errors.Is(err, validation.ErrMissingRequiredField) {
		t.Fatalf("[]string{}: %v", err)
	}
}

func TestValidateFormPayload_OptionalEmptyArrayOmitted(t *testing.T) {
	s := schema(fc.Field{
		Name:     "tags",
		Type:     fc.FieldTypeArray,
		Required: false,
		Items:    &fc.FieldItems{Type: fc.FieldTypeString},
	})

	got, err := validation.ValidateFormPayload(s, map[string]any{"tags": []any{}})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if _, ok := got["tags"]; ok {
		t.Fatalf("empty optional array should be omitted: %#v", got)
	}
}

func TestValidateFormPayload_ArrayHomogeneity(t *testing.T) {
	s := schema(fc.Field{
		Name:     "nums",
		Type:     fc.FieldTypeArray,
		Required: true,
		Items:    &fc.FieldItems{Type: fc.FieldTypeNumber},
	})
	_, err := validation.ValidateFormPayload(s, map[string]any{
		"nums": []any{int64(1), "x"},
	})
	if !errors.Is(err, validation.ErrInvalidArrayItems) {
		t.Fatalf("got %v", err)
	}
}

func TestValidateFormPayload_ArrayMissingItems(t *testing.T) {
	s := schema(fc.Field{
		Name:     "tags",
		Type:     fc.FieldTypeArray,
		Required: true,
	})
	_, err := validation.ValidateFormPayload(s, map[string]any{"tags": []any{"a"}})
	if !errors.Is(err, validation.ErrMissingArrayItemType) {
		t.Fatalf("got %v", err)
	}
}

func TestValidateFormPayload_AccumulatesErrors(t *testing.T) {
	s := schema(
		fc.Field{Name: "email", Type: fc.FieldTypeString, Required: true},
		fc.Field{Name: "age", Type: fc.FieldTypeNumber, Required: true},
	)

	_, err := validation.ValidateFormPayload(s, map[string]any{
		"age":   "nope",
		"extra": true,
	})

	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, validation.ErrMissingRequiredField) {
		t.Fatalf("missing required: %v", err)
	}

	if !errors.Is(err, validation.ErrInvalidFieldType) {
		t.Fatalf("invalid type: %v", err)
	}

	if !errors.Is(err, validation.ErrUnknownField) {
		t.Fatalf("unknown field: %v", err)
	}

	response := validation.NewValidationErrorResponse(http.StatusBadRequest, err)
	if response.Status != http.StatusBadRequest {
		t.Fatalf("status: %#v", response.Status)
	}

	if response.Message != validation.ErrValidationFailed.Error() {
		t.Fatalf("message: %#v", response.Message)
	}

	errs := response.Errors
	if len(errs) != 3 {
		t.Fatalf("errors counte: %#v", len(errs))
	}

	if errs["email"] != validation.ErrMissingRequiredField.Error() {
		t.Fatalf("email: %#v", errs["email"])
	}

	if errs["age"] != validation.ErrInvalidFieldType.Error() {
		t.Fatalf("age: %#v", errs["age"])
	}

	if errs["extra"] != validation.ErrUnknownField.Error() {
		t.Fatalf("extra: %#v", errs["extra"])
	}
}

func TestValidateFormPayload_Float64Number(t *testing.T) {
	s := schema(fc.Field{Name: "score", Type: fc.FieldTypeNumber, Required: true})
	got, err := validation.ValidateFormPayload(s, map[string]any{"score": 1.5})
	if err != nil {
		t.Fatal(err)
	}
	if got["score"] != 1.5 {
		t.Fatalf("got %#v", got)
	}
}

func TestValidateFormPayload_EmptyStringCoercionFails(t *testing.T) {
	s := schema(
		fc.Field{Name: "age", Type: fc.FieldTypeNumber, Required: true},
		fc.Field{Name: "ok", Type: fc.FieldTypeBoolean, Required: true},
	)
	_, err := validation.ValidateFormPayload(s, map[string]any{"age": "", "ok": ""})
	if !errors.Is(err, validation.ErrInvalidFieldType) {
		t.Fatalf("got %v", err)
	}
}

func TestFieldError_ErrorAndUnwrap(t *testing.T) {
	var nilErr *validation.FieldError
	if nilErr.Error() != "form field validation failed" {
		t.Fatalf("nil: %q", nilErr.Error())
	}
	if nilErr.Unwrap() != nil {
		t.Fatal("nil unwrap")
	}

	err := &validation.FieldError{Field: "email", Err: validation.ErrMissingRequiredField}
	if err.Error() != "email: "+validation.ErrMissingRequiredField.Error() {
		t.Fatalf("err = %q", err.Error())
	}
	if !errors.Is(err, validation.ErrMissingRequiredField) {
		t.Fatal("unwrap")
	}

	bare := &validation.FieldError{Err: validation.ErrUnknownField}
	if bare.Error() != validation.ErrUnknownField.Error() {
		t.Fatalf("bare = %q", bare.Error())
	}
}

func TestNewValidationErrorResponse_SingleError(t *testing.T) {
	response := validation.NewValidationErrorResponse(http.StatusBadRequest, validation.ErrUnknownField)
	if response.Message != validation.ErrUnknownField.Error() {
		t.Fatalf("message = %q", response.Message)
	}
	if response.Errors != nil {
		t.Fatalf("errors = %#v", response.Errors)
	}
}

func TestFieldErrorsMap_Empty(t *testing.T) {
	if validation.FieldErrorsMap(nil) != nil {
		t.Fatal("expected nil")
	}
}
