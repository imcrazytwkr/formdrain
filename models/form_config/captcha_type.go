package form_config

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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

func (c CaptchaType) String() string {
	return toString[c]
}

func (c CaptchaType) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(c.String())
}

func (c *CaptchaType) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	if t != bsontype.String {
		return fmt.Errorf("%q is not a valid captcha type field type", t)
	}

	var raw string
	err := bson.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	value, ok := fromString[raw]
	if !ok {
		return fmt.Errorf("%q is not a valid captcha type", raw)
	}

	*c = value
	return nil
}
