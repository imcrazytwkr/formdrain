package form_config

type FormConfig struct {
	FormId        int64           `json:"id"`
	SiteId        int64           `json:"site_id"`
	CaptchaType   CaptchaType     `json:"captcha_type"`
	CaptchaField  string          `json:"captcha_field,omitempty"`
	RedirectTo    string          `json:"redirect_to"`
	FieldSchema   FieldSchema     `json:"field_schema"`
	SchemaVersion int             `json:"schema_version"`
	Notifiers     NotifiersConfig `json:"notifiers"`
}

// CaptchaTokenField is the form map key for the captcha response token.
// CaptchaField wins when set; otherwise the provider default is used.
func (c *FormConfig) CaptchaTokenField() string {
	if c == nil {
		return ""
	}
	if len(c.CaptchaField) > 0 {
		return c.CaptchaField
	}
	return c.CaptchaType.DefaultTokenField()
}
