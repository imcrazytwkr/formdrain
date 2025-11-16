package discord

type WebhookKey struct {
	Snowflake string `bson:"snowflake"`
	Token     string `bson:"token"`
}
