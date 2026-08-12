package cursors

import (
	"encoding/base64"
	"encoding/binary"
)

func EncodeSiteIDCursor(siteID, formID int64) string {
	buf := make([]byte, 1+8+8)
	buf[0] = byte(cursorSortSiteID)
	binary.BigEndian.PutUint64(buf[1:9], uint64(siteID))
	binary.BigEndian.PutUint64(buf[9:17], uint64(formID))
	return base64.RawURLEncoding.EncodeToString(buf)
}

func DecodeSiteIDCursor(cursor string) (int64, int64, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) != 17 || cursorSort(raw[0]) != cursorSortSiteID {
		return 0, 0, errInvalidCursor
	}
	return int64(binary.BigEndian.Uint64(raw[1:9])), int64(binary.BigEndian.Uint64(raw[9:17])), nil
}
