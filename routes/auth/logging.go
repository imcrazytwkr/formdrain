package auth

import (
	"context"

	"github.com/rs/zerolog"
)

const actionForm = "form"
const actionLogin = "login"

func getLoggerForAction(ctx context.Context, action string) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("controller", "auth").Str("action", action).Logger()
}
