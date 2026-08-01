package logutil

import (
	"github.com/imcrazytwkr/formdrain/utils/errorutil"
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
	errs, ok := errorutil.UnwrapMultiErr(err)
	if ok && len(errs) > 0 {
		return event.Errs(errorsFieldName, errs)
	}

	return event.Err(err)
}
