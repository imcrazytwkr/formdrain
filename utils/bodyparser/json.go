package bodyparser

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
)

type jsonParser struct{}

func NewJsonParser() BodyParser {
	return &jsonParser{}
}

func (p *jsonParser) Parse(body []byte) (map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()

	var raw map[string]any
	if err := dec.Decode(&raw); err != nil {
		return nil, err
	}

	result := make(map[string]any, len(raw))
	for key, v := range raw {
		switch value := v.(type) {
		case []any:
			if len(value) < 1 {
				continue
			}

			parsed, err := parseArray(value)
			if err != nil {
				return nil, err
			}
			if len(parsed) > 0 {
				result[key] = parsed
			}
		default:
			parsed, err := parseSimpleValue(v)
			if err != nil {
				return nil, err
			}
			if parsed != nil {
				result[key] = parsed
			}
		}
	}

	return result, nil
}

func parseArray(arr []any) ([]any, error) {
	size := len(arr)
	if size < 1 {
		return nil, nil
	}

	var lastKind valueKind
	result := make([]any, size)
	i := 0

	for _, v := range arr {
		kind := classifyValue(v)
		switch kind {
		case kindArray, kindObject:
			return nil, errors.New("nested structures are not allowed")
		case kindNull:
			continue
		default:
			if i == 0 {
				lastKind = kind
			} else if kind != lastKind {
				return nil, errors.New("combined-type arrays are not allowed")
			}

			value, err := parseSimpleValue(v)
			if err != nil {
				return nil, err
			}
			if value != nil {
				result[i] = value
				i++
			}
		}
	}

	if i == 0 {
		return nil, nil
	}

	if i < size {
		return result[:i], nil
	}

	return result, nil
}

type valueKind int8

const (
	kindNull valueKind = iota
	kindString
	kindNumber
	kindBool
	kindArray
	kindObject
	kindOther
)

func classifyValue(v any) valueKind {
	switch v.(type) {
	case nil:
		return kindNull
	case string:
		return kindString
	case json.Number:
		return kindNumber
	case bool:
		return kindBool
	case []any:
		return kindArray
	case map[string]any:
		return kindObject
	default:
		return kindOther
	}
}

func parseSimpleValue(v any) (any, error) {
	switch value := v.(type) {
	case nil:
		return nil, nil
	case string:
		return strings.TrimSpace(value), nil
	case json.Number:
		return parseNumber(value)
	case bool:
		return value, nil
	case map[string]any:
		// Top-level nested objects are skipped (same as previous fastjson path).
		return nil, nil
	default:
		return nil, nil
	}
}

func parseNumber(n json.Number) (any, error) {
	intValue, err := n.Int64()
	if err == nil {
		return intValue, nil
	}

	return n.Float64()
}
