package errorutil

// joinErrors
type joinError interface {
	Unwrap() []error
}

/**
 * Unsraps a joinError. Will **not** work well with `wrapErrors`
 * that has a dedicated message.
 */
func UnwrapMultiErr(err error) ([]error, bool) {
	e, ok := err.(joinError)
	if !ok {
		return nil, false
	}

	return e.Unwrap(), true
}
