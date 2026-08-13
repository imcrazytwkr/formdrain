package captcha_validation

import (
	"fmt"

	"github.com/imcrazytwkr/formdrain/models/common"
)

func getErrNoCaptchaImpl(captchaType common.CaptchaType) error {
	return fmt.Errorf("could not find any implementation for catcha type %s", captchaType)
}

func getErrMissingCaptchaSecret(captchaType common.CaptchaType) error {
	return fmt.Errorf("missing captcha secret for type %s", captchaType)
}
