package discord

type Author struct {
	Name string `json:"name"`
	Url  string `json:"url,omitempty"`
	Icon string `json:"icon_url,omitempty"`
}
