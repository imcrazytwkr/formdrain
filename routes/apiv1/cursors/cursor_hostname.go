package cursors

import (
	"encoding/base64"
	"encoding/binary"
)

func EncodeHostnameCursor(hostname string, id int64) string {
	host := []byte(hostname)
	buf := make([]byte, 1+8+len(host))
	buf[0] = byte(cursorSortHostname)
	binary.BigEndian.PutUint64(buf[1:9], uint64(id))
	copy(buf[9:], host)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func DecodeHostnameCursor(cursor string) (string, int64, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) < 9 || cursorSort(raw[0]) != cursorSortHostname {
		return "", 0, errInvalidCursor
	}

	return string(raw[9:]), int64(binary.BigEndian.Uint64(raw[1:9])), nil
}
