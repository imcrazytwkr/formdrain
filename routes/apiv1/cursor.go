package apiv1

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"time"

	"github.com/google/uuid"
)

type cursorSort byte

const (
	cursorSortDefault cursorSort = iota
	cursorSortID
	cursorSortHostname
	cursorSortSiteID
	cursorSortCreatedAt
)

var errInvalidCursor = errors.New("invalid cursor")

func encodeIDCursor(id int64) string {
	buf := make([]byte, 1+8)
	buf[0] = byte(cursorSortID)
	binary.BigEndian.PutUint64(buf[1:], uint64(id))
	return base64.RawURLEncoding.EncodeToString(buf)
}

func encodeHostnameCursor(hostname string, id int64) string {
	host := []byte(hostname)
	buf := make([]byte, 1+8+len(host))
	buf[0] = byte(cursorSortHostname)
	binary.BigEndian.PutUint64(buf[1:9], uint64(id))
	copy(buf[9:], host)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func decodeIDCursor(cursor string) (int64, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) != 9 || cursorSort(raw[0]) != cursorSortID {
		return 0, errInvalidCursor
	}

	return int64(binary.BigEndian.Uint64(raw[1:])), nil
}

func decodeHostnameCursor(cursor string) (string, int64, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) < 9 || cursorSort(raw[0]) != cursorSortHostname {
		return "", 0, errInvalidCursor
	}

	return string(raw[9:]), int64(binary.BigEndian.Uint64(raw[1:9])), nil
}

func encodeSiteIDCursor(siteID, formID int64) string {
	buf := make([]byte, 1+8+8)
	buf[0] = byte(cursorSortSiteID)
	binary.BigEndian.PutUint64(buf[1:9], uint64(siteID))
	binary.BigEndian.PutUint64(buf[9:17], uint64(formID))
	return base64.RawURLEncoding.EncodeToString(buf)
}

func decodeSiteIDCursor(cursor string) (int64, int64, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) != 17 || cursorSort(raw[0]) != cursorSortSiteID {
		return 0, 0, errInvalidCursor
	}
	return int64(binary.BigEndian.Uint64(raw[1:9])), int64(binary.BigEndian.Uint64(raw[9:17])), nil
}

// We're currently using UUID v4
const idLength = 36

func encodeCreatedAtCursor(createdAt time.Time, id string) string {
	buf := make([]byte, 1+8+idLength)
	buf[0] = byte(cursorSortCreatedAt)
	binary.BigEndian.PutUint64(buf[1:9], uint64(createdAt.UTC().Unix()))
	copy(buf[9:], id)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func decodeCreatedAtCursor(cursor string) (time.Time, string, error) {
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
