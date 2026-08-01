package recaptcha

type recaptchaResponse struct {
	Success  bool   `json:"success"`
	Hostname string `json:"hostname"`
}
