package cursors

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"testing"
)

func TestSiteIDCursor_RoundTrip(t *testing.T) {
	t.Parallel()

	cursor := EncodeSiteIDCursor(3, 99)
	siteID, formID, err := DecodeSiteIDCursor(cursor)
	if err != nil {
		t.Fatal(err)
	}
	if siteID != 3 || formID != 99 {
		t.Fatalf("got site=%d form=%d", siteID, formID)
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
		EncodeIDCursor(1),
		EncodeHostnameCursor("x", 1),
		base64.RawURLEncoding.EncodeToString(wrongLen),
	}

	for _, cursor := range cases {
		_, _, err := DecodeSiteIDCursor(cursor)
		if !errors.Is(err, errInvalidCursor) {
			t.Fatalf("cursor %q: err=%v", cursor, err)
		}
	}
}
