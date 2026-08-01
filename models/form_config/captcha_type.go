package form_config

import (
	"encoding/json"
	"fmt"
)

type CaptchaType int8

const (
	CaptchaTypeUndefined CaptchaType = iota
	CaptchaTypeHcaptcha
	CaptchaTypeRecaptcha
)

var toString = map[CaptchaType]string{
	CaptchaTypeUndefined: "undefined",
	CaptchaTypeHcaptcha:  "hcaptcha",
	CaptchaTypeRecaptcha: "recaptcha",
}

var fromString map[string]CaptchaType

func init() {
	fromString = make(map[string]CaptchaType, len(toString))
	for k, v := range toString {
		fromString[v] = k
	}
}

func ParseCaptchaType(v string) (CaptchaType, error) {
	value, ok := fromString[v]
	if !ok {
		return CaptchaTypeUndefined, fmt.Errorf("%q is not a valid captcha type", v)
	}

	return value, nil
}

func (c CaptchaType) String() string {
	return toString[c]
}

func (c CaptchaType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *CaptchaType) UnmarshalJSON(b []byte) error {
	var raw string
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	value, err := ParseCaptchaType(raw)
	if err != nil {
		return err
	}

	*c = value
	return nil
}
