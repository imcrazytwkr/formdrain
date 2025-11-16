package models

type Request struct {
	Username string   `json:"username"`
	Avatar   string   `json:"avatar_url"`
	Content  string   `json:"content,omitempty"`
	Embeds   []*Embed `json:"embeds,omitempty"`
}
