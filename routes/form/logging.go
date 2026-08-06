package form

import (
	"context"

	"github.com/rs/zerolog"
)

const actionSend = "send"

func getLoggerForAction(ctx context.Context, action string) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("controller", "form").Str("action", action).Logger()
}
