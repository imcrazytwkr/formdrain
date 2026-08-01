package bodyparser

import "errors"

var ErrNestedStructuresNotAllowed = errors.New("nested structures are not allowed")
var ErrCombinedTypeArraysNotAllowed = errors.New("combined-type arrays are not allowed")
