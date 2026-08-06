package form

import (
	"errors"
	"fmt"
)

var errInvalidOrigin = errors.New("CORS origin is not valid")

var errFormTooLarge = errors.New("forms data can only be up to 4kb in size")

func getErrFormNotFound(id string) error {
	return fmt.Errorf("could not find form with ID %q", id)
}
