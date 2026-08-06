package bodyparser

import "errors"

var ErrNestedStructuresNotAllowed = errors.New("nested structures are not allowed")
var ErrCombinedTypeArraysNotAllowed = errors.New("combined-type arrays are not allowed")

// "forms data can only be up to 4kb in size"
var ErrUnsupportedContentType = errors.New("unsupported content type")
var ErrBodyTooLarge = errors.New("request body can be at most 4kb")
var ErrMalformedBody = errors.New("malformed body")
