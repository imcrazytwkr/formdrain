package form_config

import "errors"

var ErrInvalidCaptchaType = errors.New(`Captcha type "undefined" is not allowed in DB`)
