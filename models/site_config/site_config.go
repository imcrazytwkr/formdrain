package site_config

import (
	"github.com/imcrazytwkr/formdrain/models/common"
)

type SiteConfig struct {
	SiteId   common.ObjectId `bson:"_id"`
	Hostname string          `bson:"hostname"`
}
