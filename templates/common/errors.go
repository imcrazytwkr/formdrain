package common

import "errors"

var ErrNilTemplate = errors.New("common: nil template")
var ErrTooLarge = errors.New("common: template too large")
var ErrEmpty = errors.New("common: empty template")
