package safemustache

import "errors"

var ErrRawInterpolation = errors.New("safemustache: unescaped interpolation is not allowed")
var ErrDelimiterChange = errors.New("safemustache: changing delimiters is not allowed")
var ErrPartial = errors.New("safemustache: partials are not allowed")
