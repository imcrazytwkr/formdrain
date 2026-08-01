package logutil

import (
	"github.com/rs/zerolog"
)

const errorsFieldName = "errors"

// joinErrors
type joinError interface {
	Unwrap() []error
}

/**
 * Unsraps a joinError. Will **not** work well with `wrapErrors`
 * that has a dedicated message.
 */
func UnwrapErr(event *zerolog.Event, err error) *zerolog.Event {
	e, ok := err.(joinError)
	if ok {
		return event.Errs(errorsFieldName, e.Unwrap())
	}

	return event.Err(err)
}
