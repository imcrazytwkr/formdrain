package cursors

import (
	"encoding/base64"
	"encoding/binary"
)

func EncodeIDCursor(id int64) string {
	buf := make([]byte, 1+8)
	buf[0] = byte(cursorSortID)
	binary.BigEndian.PutUint64(buf[1:], uint64(id))
	return base64.RawURLEncoding.EncodeToString(buf)
}

func DecodeIDCursor(cursor string) (int64, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) != 9 || cursorSort(raw[0]) != cursorSortID {
		return 0, errInvalidCursor
	}

	return int64(binary.BigEndian.Uint64(raw[1:])), nil
}
