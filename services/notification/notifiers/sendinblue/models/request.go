package models

import "github.com/imcrazytwkr/formdrain/models/form_config/sendinblue"

type Request struct {
	Sender     *sendinblue.EmailContact   `json:"sender"`
	Recipients []*sendinblue.EmailContact `json:"to"`
	Subject    string                     `json:"subject"`
	Content    string                     `json:"htmlContent"`
}
