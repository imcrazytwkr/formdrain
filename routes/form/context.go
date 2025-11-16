package form

import (
	"github.com/gin-gonic/gin"
	"github.com/imcrazytwkr/formdrain/models/form_config"
	"github.com/imcrazytwkr/formdrain/models/site_config"
)

func getFormConfig(c *gin.Context) (*form_config.FormConfig, bool) {
	rawFormConfig, ok := c.Get(keyFormConfig)
	if !ok {
		return nil, false
	}

	formConfig, ok := rawFormConfig.(*form_config.FormConfig)
	return formConfig, ok
}

func getSiteConfig(c *gin.Context) (*site_config.SiteConfig, bool) {
	rawSiteConfig, ok := c.Get(keySiteConfig)
	if !ok {
		return nil, false
	}

	siteConfig, ok := rawSiteConfig.(*site_config.SiteConfig)
	return siteConfig, ok
}
