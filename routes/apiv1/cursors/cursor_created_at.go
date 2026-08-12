package cursors

import (
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/google/uuid"
)

// We're currently using UUID v4
const idLength = 36

func EncodeCreatedAtCursor(createdAt time.Time, id string) string {
	buf := make([]byte, 1+8+idLength)
	buf[0] = byte(cursorSortCreatedAt)
	binary.BigEndian.PutUint64(buf[1:9], uint64(createdAt.UTC().Unix()))
	copy(buf[9:], id)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func DecodeCreatedAtCursor(cursor string) (time.Time, string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) != 9+idLength || cursorSort(raw[0]) != cursorSortCreatedAt {
		return time.Time{}, "", errInvalidCursor
	}

	id := string(raw[9:])
	if uuid.Validate(id) != nil {
		return time.Time{}, "", errInvalidCursor
	}

	sec := int64(binary.BigEndian.Uint64(raw[1:9]))
	return time.Unix(sec, 0).UTC(), id, nil
}
