package bodyparser

import (
	"bytes"
	"errors"

	"github.com/valyala/fastjson"
)

type jsonParser struct {
	pool *fastjson.ParserPool
}

func NewJsonParser() BodyParser {
	pool := &fastjson.ParserPool{}
	return &jsonParser{pool}
}

func (p *jsonParser) Parse(body []byte) (map[string]any, error) {
	var err error

	parser := p.pool.Get()
	defer p.pool.Put(parser)

	bodyValue, err := parser.ParseBytes(body)
	if err != nil {
		return nil, err
	}

	object, err := bodyValue.Object()
	if err != nil {
		return nil, err
	}

	var result map[string]any
	object.Visit(func(k []byte, v *fastjson.Value) {
		if err != nil {
			return
		}

		key := string(k)
		switch v.Type() {
		case fastjson.TypeArray:
			values, e := v.Array()
			if e != nil || len(values) < 1 {
				break
			}

			// @NOTE: error caught here means request body is FUBAR
			value, e := parseArray(values)
			if e != nil {
				err = e
				break
			}

			if len(value) > 0 {
				result[key] = value
			}
		default:
			value, e := parseSimpleValue(v)
			if e == nil && value != nil {
				result[key] = value
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseArray(arr []*fastjson.Value) ([]any, error) {
	size := len(arr)
	if size < 1 {
		return nil, nil
	}

	lastType := fastjson.TypeNull
	result := make([]any, size)
	i := 0

	for _, v := range arr {
		switch v.Type() {
		case fastjson.TypeArray, fastjson.TypeObject:
			return nil, errors.New("nested structures are not allowed")
		case fastjson.TypeNull:
			continue
		default:
			if i == 0 {
				lastType = v.Type()
			} else if v.Type() != lastType {
				return nil, errors.New("combined-type arrays are not allowed")
			}

			value, err := parseSimpleValue(v)
			if err == nil && value != nil {
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

func parseSimpleValue(v *fastjson.Value) (value any, err error) {
	switch v.Type() {
	case fastjson.TypeString:
		value, err = parseString(v)
	case fastjson.TypeNumber:
		value, err = parseNumber(v)
	case fastjson.TypeTrue, fastjson.TypeFalse:
		value, err = v.Bool()
	}

	return
}

func parseString(v *fastjson.Value) (string, error) {
	raw, err := v.StringBytes()
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(raw)), nil
}

func parseNumber(v *fastjson.Value) (any, error) {
	floatValue, err := v.Float64()
	if err != nil {
		return nil, err
	}

	intValue, err := v.Int64()
	if err != nil {
		return floatValue, nil
	}

	if float64(intValue) != floatValue {
		return floatValue, nil
	}

	return intValue, nil
}
