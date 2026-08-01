package form_config

type FormConfig struct {
	FormId        int64            `json:"id"`
	SiteId        int64            `json:"site_id"`
	CaptchaType   CaptchaType      `json:"captcha_type"`
	RedirectTo    string           `json:"redirect_to"`
	FieldSchema   FieldSchema      `json:"field_schema"`
	SchemaVersion int              `json:"schema_version"`
	Notifiers     *NotifiersConfig `json:"notifiers"`
}
