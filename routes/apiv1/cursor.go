package apiv1

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
)

type cursorSort byte

const (
	cursorSortDefault cursorSort = iota
	cursorSortID
	cursorSortHostname
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
