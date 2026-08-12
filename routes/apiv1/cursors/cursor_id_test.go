package cursors

import (
	"encoding/base64"
	"errors"
	"testing"
)

func TestIDCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	cursor := EncodeIDCursor(42)
	got, err := DecodeIDCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("got %d", got)
	}
}

func TestDecodeIDCursor_Invalid(t *testing.T) {
	t.Parallel()

	cases := []string{
		"not-valid",
		"",
		EncodeHostnameCursor("x", 1),
		EncodeSiteIDCursor(1, 2),
		base64.RawURLEncoding.EncodeToString([]byte{byte(cursorSortID)}),
	}

	for _, cursor := range cases {
		_, err := DecodeIDCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}
