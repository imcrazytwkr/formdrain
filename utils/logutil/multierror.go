package logutil

import (
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
)

func Multierr(event *zerolog.Event, err error) *zerolog.Event {
	errs, ok := err.(*multierror.Error)
	if ok {
		return event.Errs("errors", errs.Errors)
	}

	return event.Err(err)
}
