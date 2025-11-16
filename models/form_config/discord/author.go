package discord

type Author struct {
	Name string `bson:"name" json:"name"`
	Url  string `bson:"url,omitempty" json:"url,omitempty"`
	Icon string `bson:"icon_url,omitempty" json:"icon_url,omitempty"`
}
