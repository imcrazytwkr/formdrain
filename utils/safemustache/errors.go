package safemustache

import "errors"

var ErrTooLarge = errors.New("safemustache: template too large")
var ErrEmpty = errors.New("safemustache: empty template")
var ErrRawInterpolation = errors.New("safemustache: unescaped interpolation is not allowed")
var ErrDelimiterChange = errors.New("safemustache: changing delimiters is not allowed")
var ErrPartial = errors.New("safemustache: partials are not allowed")
