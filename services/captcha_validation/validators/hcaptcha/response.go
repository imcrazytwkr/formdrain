package hcaptcha

type hcaptchaResponse struct {
	Success  bool   `json:"success"`
	Hostname string `json:"hostname"`
}
