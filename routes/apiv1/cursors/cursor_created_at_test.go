package cursors

import (
	"errors"
	"testing"
	"time"
)

func TestCreatedAtCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	id := "00000000-0000-4000-8000-00000000000a"
	cursor := EncodeCreatedAtCursor(created, id)
	gotTime, gotID, err := DecodeCreatedAtCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if !gotTime.Equal(created) || gotID != id {
		t.Fatalf("got time=%s id=%q", gotTime, gotID)
	}
}

func TestDecodeCreatedAtCursor_Invalid(t *testing.T) {
	t.Parallel()

	cases := []string{
		"not-valid",
		"",
		EncodeIDCursor(1),
		EncodeHostnameCursor("x", 1),
		EncodeSiteIDCursor(1, 2),
		EncodeCreatedAtCursor(time.Unix(1, 0).UTC(), "not-a-uuid"),
	}

	for _, cursor := range cases {
		_, _, err := DecodeCreatedAtCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}
