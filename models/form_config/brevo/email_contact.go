package brevo

type EmailContact struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"email"`
}
