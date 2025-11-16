package models

import "github.com/imcrazytwkr/formdrain/models/form_config/discord"

type Embed struct {
	Author      *discord.Author `json:"author"`
	Title       string          `json:"title,omitempty"`
	Url         string          `json:"url,omitempty"`
	Description string          `json:"description,omitempty"`
	Color       int             `json:"color"`
	Fields      []*Field        `json:"fields,omitempty"`
	Thumbnail   *Image          `json:"thumbnail,omitempty"`
	Image       *Image          `json:"image,omitempty"`
	Footer      *Footer         `json:"footer,omitempty"`
}
