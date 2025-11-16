package bodyparser

import (
	"net/url"
	"strings"
)

type formParser struct{}

func NewFormParser() BodyParser {
	return &formParser{}
}

func (p *formParser) Parse(body []byte) (map[string]any, error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	result := make(map[string]any, len(values))
	for k, v := range values {
		if len(v) == 1 {
			value := cleanString(v[0])
			if len(value) > 0 {
				result[k] = value
			}

			continue
		}

		value := cleanSlice(v)
		if len(value) > 0 {
			result[k] = value
		}
	}

	return result, nil
}

func cleanString(s string) string {
	value := strings.TrimSpace(s)
	if value == "null" {
		return ""
	}

	return value
}

func cleanSlice(s []string) []string {
	result := make([]string, len(s))
	i := 0

	for _, v := range s {
		value := cleanString(v)
		if len(value) > 0 {
			result[i] = value
			i++
		}
	}

	if i < len(s) {
		return result[:i]
	}

	return result
}
