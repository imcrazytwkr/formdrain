package sendinblue

import "github.com/imcrazytwkr/formdrain/models/common"

type SendinblueConfig struct {
	Recipients []*EmailContact  `bson:"recipients"`
	Subject    string           `bson:"subject"`
	Template   *common.Template `bson:"template"`
}

func (c *SendinblueConfig) RenderContent(form map[string]any) string {
	return c.Template.ExecuteString(form)
}
