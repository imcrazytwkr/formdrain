package captcha_validation

import (
	"fmt"

	fc "github.com/imcrazytwkr/formdrain/models/form_config"
)

func getErrNoCaptchaImpl(captchaType fc.CaptchaType) error {
	return fmt.Errorf("could not find any implementation for catcha type %s", captchaType)
}
