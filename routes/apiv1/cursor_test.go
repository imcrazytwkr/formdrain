package apiv1

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"testing"
	"time"
)

func TestIDCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	cursor := encodeIDCursor(42)
	got, err := decodeIDCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if got != 42 {
		t.Fatalf("got %d", got)
	}
}

func TestHostnameCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	cursor := encodeHostnameCursor("a.example", 7)
	host, id, err := decodeHostnameCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if host != "a.example" || id != 7 {
		t.Fatalf("got host=%q id=%d", host, id)
	}

	cursor = encodeHostnameCursor("", 1)
	host, id, err = decodeHostnameCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if host != "" || id != 1 {
		t.Fatalf("empty host: host=%q id=%d", host, id)
	}
}

func TestSiteIDCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	cursor := encodeSiteIDCursor(3, 99)
	siteID, formID, err := decodeSiteIDCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if siteID != 3 || formID != 99 {
		t.Fatalf("got site=%d form=%d", siteID, formID)
	}
}

func TestDecodeIDCursor_Invalid(t *testing.T) {
	t.Parallel()

	cases := []string{
		"not-valid",
		"",
		encodeHostnameCursor("x", 1),
		encodeSiteIDCursor(1, 2),
		base64.RawURLEncoding.EncodeToString([]byte{byte(cursorSortID)}), // too short
	}

	for _, cursor := range cases {
		_, err := decodeIDCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}

func TestDecodeHostnameCursor_Invalid(t *testing.T) {
	t.Parallel()

	cases := []string{
		"not-valid",
		"",
		encodeIDCursor(1),
		encodeSiteIDCursor(1, 2),
		base64.RawURLEncoding.EncodeToString([]byte{byte(cursorSortHostname), 1, 2, 3}), // too short
	}

	for _, cursor := range cases {
		_, _, err := decodeHostnameCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}

func TestDecodeSiteIDCursor_Invalid(t *testing.T) {
	t.Parallel()

	wrongLen := make([]byte, 9)
	wrongLen[0] = byte(cursorSortSiteID)
	binary.BigEndian.PutUint64(wrongLen[1:], 1)

	cases := []string{
		"not-valid",
		"",
		encodeIDCursor(1),
		encodeHostnameCursor("x", 1),
		base64.RawURLEncoding.EncodeToString(wrongLen),
	}

	for _, cursor := range cases {
		_, _, err := decodeSiteIDCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}

func TestCreatedAtCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	created := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	id := "00000000-0000-4000-8000-00000000000a"
	cursor := encodeCreatedAtCursor(created, id)
	gotTime, gotID, err := decodeCreatedAtCursor(cursor)
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
		encodeIDCursor(1),
		encodeHostnameCursor("x", 1),
		encodeSiteIDCursor(1, 2),
		encodeCreatedAtCursor(time.Unix(1, 0).UTC(), "not-a-uuid"),
	}

	for _, cursor := range cases {
		_, _, err := decodeCreatedAtCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}
