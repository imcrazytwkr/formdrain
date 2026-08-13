package site_config

import (
	"github.com/imcrazytwkr/formdrain/models/common"
)

type SiteConfig struct {
	SiteId          int64  `json:"id"`
	Hostname        string `json:"hostname"`
	OwnerId         int64  `json:"owner_id"`
	HcaptchaSecret  string `json:"-"`
	RecaptchaSecret string `json:"-"`
}

func (c *SiteConfig) CaptchaSecret(t common.CaptchaType) string {
	if c == nil {
		return ""
	}
	switch t {
	case common.CaptchaTypeHcaptcha:
		return c.HcaptchaSecret
	case common.CaptchaTypeRecaptcha:
		return c.RecaptchaSecret
	default:
		return ""
	}
}
