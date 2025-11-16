package form_config

import (
	"github.com/imcrazytwkr/formdrain/models/common"
)

type FormConfig struct {
	FormId      common.ObjectId  `bson:"_id"`
	SiteId      common.ObjectId  `bson:"site_id"`
	CaptchaType CaptchaType      `bson:"captcha_type"`
	Notifiers   *NotifiersConfig `bson:"notifiers"`
	RedirectTo  string           `bson:"redirect_to"`
}
