package site_config

type SiteConfig struct {
	SiteId   int64  `json:"id"`
	Hostname string `json:"hostname"`
	OwnerId  int64  `json:"owner_id"`
}
