package common

import (
	"context"

	"github.com/rs/zerolog"
)

const ApiFormatHttp = "http"

func GetLoggerForProvider(ctx context.Context, provider string, format string) zerolog.Logger {
	return zerolog.Ctx(ctx).With().
		Str("service", "captcha_validation").
		Str("provider", provider).
		Str("format", format).
		Logger()
}
