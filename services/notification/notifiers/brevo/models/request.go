package models

import "github.com/imcrazytwkr/formdrain/models/form_config/brevo"

type Request struct {
	Sender     *brevo.EmailContact   `json:"sender"`
	Recipients []*brevo.EmailContact `json:"to"`
	Subject    string                `json:"subject"`
	Content    string                `json:"htmlContent"`
}
