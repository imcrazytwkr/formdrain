package sendinblue

import "github.com/imcrazytwkr/formdrain/models/common"

type SendinblueConfig struct {
	Recipients []*EmailContact  `json:"recipients"`
	Subject    string           `json:"subject"`
	Template   *common.Template `json:"template"`
}

func (c *SendinblueConfig) RenderContent(form map[string]any) string {
	return c.Template.ExecuteString(form)
}
