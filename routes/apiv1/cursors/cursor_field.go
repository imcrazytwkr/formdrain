package cursors

import (
	"encoding/base64"
	"encoding/binary"
	"math"

	"github.com/google/uuid"
	fc "github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/utils/collate"
)

const (
	cursorFieldString byte = iota + 1
	cursorFieldNumber
	cursorFieldBoolean
)

const (
	cursorFieldFlagDesc byte = 1 << 0
	cursorFieldFlagNull byte = 1 << 1
)

type FieldCursor struct {
	Field string
	Type  fc.FieldType
	Desc  bool
	Null  bool
	Value any
	ID    string
}

func EncodeFieldCursor(c FieldCursor) (string, error) {
	typeByte, ok := cursorFieldType(c.Type)
	if !ok || len(c.Field) < 1 || len(c.Field) > 0xffff || uuid.Validate(c.ID) != nil {
		return "", errInvalidCursor
	}

	var flags byte
	if c.Desc {
		flags |= cursorFieldFlagDesc
	}
	if c.Null {
		flags |= cursorFieldFlagNull
	}

	value, err := encodeFieldCursorValue(c)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 1+1+1+idLength+2+len(c.Field)+len(value))
	buf[0] = byte(cursorSortField)
	buf[1] = typeByte
	buf[2] = flags
	copy(buf[3:3+idLength], c.ID)
	binary.BigEndian.PutUint16(buf[3+idLength:5+idLength], uint16(len(c.Field)))
	copy(buf[5+idLength:], c.Field)
	copy(buf[5+idLength+len(c.Field):], value)
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func encodeFieldCursorValue(c FieldCursor) ([]byte, error) {
	if c.Null {
		return nil, nil
	}

	switch c.Type {
	case fc.FieldTypeString:
		s, ok := c.Value.(string)
		if !ok {
			return nil, errInvalidCursor
		}
		return []byte(s), nil
	case fc.FieldTypeNumber:
		f, ok := collate.NumberToFloat64(c.Value)
		if !ok {
			return nil, errInvalidCursor
		}

		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, math.Float64bits(f))
		return buf, nil
	case fc.FieldTypeBoolean:
		b, ok := c.Value.(bool)
		if !ok {
			return nil, errInvalidCursor
		}

		if b {
			return []byte{1}, nil
		}

		return []byte{0}, nil
	default:
		return nil, errInvalidCursor
	}
}

func cursorFieldType(fieldType fc.FieldType) (byte, bool) {
	switch fieldType {
	case fc.FieldTypeString:
		return cursorFieldString, true
	case fc.FieldTypeNumber:
		return cursorFieldNumber, true
	case fc.FieldTypeBoolean:
		return cursorFieldBoolean, true
	default:
		return 0, false
	}
}

func FloatSortValue(value any) (float64, bool) {
	switch v := value.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	default:
		return 0, false
	}
}

func DecodeFieldCursor(cursor string) (FieldCursor, error) {
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) < 5+idLength || cursorSort(raw[0]) != cursorSortField {
		return FieldCursor{}, errInvalidCursor
	}

	fieldType, ok := fieldTypeFromCursor(raw[1])
	if !ok {
		return FieldCursor{}, errInvalidCursor
	}

	flags := raw[2]
	id := string(raw[3 : 3+idLength])
	if uuid.Validate(id) != nil {
		return FieldCursor{}, errInvalidCursor
	}

	nameLen := int(binary.BigEndian.Uint16(raw[3+idLength : 5+idLength]))
	nameStart := 5 + idLength
	if nameLen < 1 || nameStart+nameLen > len(raw) {
		return FieldCursor{}, errInvalidCursor
	}

	c := FieldCursor{
		Field: string(raw[nameStart : nameStart+nameLen]),
		Type:  fieldType,
		Desc:  flags&cursorFieldFlagDesc != 0,
		Null:  flags&cursorFieldFlagNull != 0,
		ID:    id,
	}

	value := raw[nameStart+nameLen:]
	if c.Null {
		if len(value) != 0 {
			return FieldCursor{}, errInvalidCursor
		}
		return c, nil
	}

	c.Value, err = decodeFieldCursorValue(fieldType, value)
	if err != nil {
		return FieldCursor{}, err
	}
	return c, nil
}

func decodeFieldCursorValue(fieldType fc.FieldType, raw []byte) (any, error) {
	switch fieldType {
	case fc.FieldTypeString:
		return string(raw), nil
	case fc.FieldTypeNumber:
		if len(raw) != 8 {
			return nil, errInvalidCursor
		}
		return math.Float64frombits(binary.BigEndian.Uint64(raw)), nil
	case fc.FieldTypeBoolean:
		if len(raw) != 1 || (raw[0] != 0 && raw[0] != 1) {
			return nil, errInvalidCursor
		}
		return raw[0] == 1, nil
	default:
		return nil, errInvalidCursor
	}
}

func fieldTypeFromCursor(b byte) (fc.FieldType, bool) {
	switch b {
	case cursorFieldString:
		return fc.FieldTypeString, true
	case cursorFieldNumber:
		return fc.FieldTypeNumber, true
	case cursorFieldBoolean:
		return fc.FieldTypeBoolean, true
	default:
		return "", false
	}
}
