package cursors

import (
	"errors"
	"testing"
	"time"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

func TestFieldCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	id := "00000000-0000-4000-8000-00000000000a"
	cases := []FieldCursor{
		{Field: "email", Type: fc.FieldTypeString, ID: id, Value: "a"},
		{Field: "email", Type: fc.FieldTypeString, ID: id, Value: ""},
		{Field: "email", Type: fc.FieldTypeString, ID: id, Null: true},
		{Field: "age", Type: fc.FieldTypeNumber, Desc: true, ID: id, Value: 2.5},
		{Field: "ok", Type: fc.FieldTypeBoolean, ID: id, Value: true},
		{Field: "ok", Type: fc.FieldTypeBoolean, ID: id, Value: false},
		{Field: `user.email`, Type: fc.FieldTypeString, ID: id, Value: "z"},
	}

	for _, want := range cases {
		cursor, err := EncodeFieldCursor(want)
		if err != nil {
			t.Fatalf("%#v: encode: %v", want, err)
		}
		got, err := DecodeFieldCursor(cursor)
		if err != nil {
			t.Fatalf("%#v: decode: %v", want, err)
		}
		if got.Field != want.Field || got.Type != want.Type || got.Desc != want.Desc || got.Null != want.Null || got.ID != want.ID {
			t.Fatalf("got %#v want %#v", got, want)
		}
		if !want.Null && got.Value != want.Value {
			t.Fatalf("value got %#v want %#v", got.Value, want.Value)
		}
	}
}

func TestDecodeFieldCursor_Invalid(t *testing.T) {
	t.Parallel()

	id := "00000000-0000-4000-8000-00000000000a"
	cases := []string{
		"not-valid",
		"",
		EncodeIDCursor(1),
		EncodeCreatedAtCursor(time.Unix(1, 0).UTC(), id),
	}

	for _, cursor := range cases {
		_, err := DecodeFieldCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}
