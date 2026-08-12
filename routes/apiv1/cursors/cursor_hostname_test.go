package cursors

import (
	"encoding/base64"
	"errors"
	"testing"
)

func TestHostnameCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	cursor := EncodeHostnameCursor("a.example", 7)
	host, id, err := DecodeHostnameCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if host != "a.example" || id != 7 {
		t.Fatalf("got host=%q id=%d", host, id)
	}

	cursor = EncodeHostnameCursor("", 1)
	host, id, err = DecodeHostnameCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if host != "" || id != 1 {
		t.Fatalf("empty host: host=%q id=%d", host, id)
	}
}

func TestDecodeHostnameCursor_Invalid(t *testing.T) {
	t.Parallel()

	cases := []string{
		"not-valid",
		"",
		EncodeIDCursor(1),
		EncodeSiteIDCursor(1, 2),
		base64.RawURLEncoding.EncodeToString([]byte{byte(cursorSortHostname), 1, 2, 3}),
	}

	for _, cursor := range cases {
		_, _, err := DecodeHostnameCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}
